<template>
  <main-layout>
    <el-card class="form-card">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="150px" class="custom-form">
        <el-form-item :label="t('form.senderName')" prop="name">
          <el-input
            v-model="form.name"
            maxlength="100"
            :placeholder="t('form.namePlaceholder')"
            class="input-field"
          />
        </el-form-item>

        <el-form-item :label="t('form.deliveryDate')" prop="date">
          <el-date-picker
            v-model="form.date"
            type="date"
            format="DD/MM/YYYY"
            :placeholder="t('form.pickDate')"
            :disabled-date="isNotFutureDate"
            class="input-field date-field"
          />
        </el-form-item>

        <el-form-item :label="t('form.message')" prop="message">
          <el-input
            v-model="form.message"
            type="textarea"
            maxlength="4096"
            show-word-limit
            :placeholder="t('form.message')"
            class="input-field custom-input"
          />
        </el-form-item>

        <el-form-item :label="t('form.recipientEmail')" prop="email">
          <el-input
            v-model="form.email"
            maxlength="255"
            :placeholder="t('form.recipientEmail')"
            class="input-field"
          />
        </el-form-item>

        <el-form-item :label="t('form.attachments')">
          <div class="attachment-container">
            <el-upload
              class="file-upload"
              list-type="picture-card"
              accept="image/*"
              :limit="MAX_ATTACHMENTS"
              :http-request="uploadAttachmentRequest"
              :before-upload="validateAttachment"
              :on-remove="handleRemove"
              :on-exceed="handleExceed"
            >
              <el-icon><Plus /></el-icon>
            </el-upload>
          </div>
        </el-form-item>

        <el-form-item class="form-buttons">
          <div class="submit-reset-container">
            <el-button
              type="primary"
              class="submit-button"
              :loading="isSubmitting"
              @click="submitForm"
            >
              {{ t('form.submit') }}
            </el-button>
            <el-button type="default" class="reset-button" @click="resetForm">
              {{ t('form.reset') }}
            </el-button>
          </div>
        </el-form-item>
      </el-form>
    </el-card>
  </main-layout>
</template>

<script setup>
import { computed, reactive, ref } from 'vue';
import { v4 as uuidv4 } from 'uuid';
import dayjs from 'dayjs';
import { fileTypeFromBuffer } from 'file-type';
import { ElMessage } from 'element-plus';
import { Plus } from '@element-plus/icons-vue';
import { useI18n } from 'vue-i18n';
import MainLayout from '@/layouts/MainLayout.vue';
import { createCapsule, getUploadTarget, uploadAttachment } from '@/api/capsules';

// Mirrors the signed upload policy. Checking here only buys faster feedback.
const MAX_ATTACHMENTS = 3;
const MAX_ATTACHMENT_BYTES = 5 * 1024 * 1024;
const ALLOWED_IMAGE_TYPES = ['image/jpeg', 'image/png', 'image/webp', 'image/gif'];

const { t } = useI18n();

const formRef = ref(null);
const isSubmitting = ref(false);
const uploadedCount = ref(0);
// One folder per capsule links its uploads to the row, without the backend ever
// seeing the files.
const filesFolderUuid = ref(uuidv4());

const form = reactive({ name: '', date: '', message: '', email: '' });

const rules = computed(() => ({
  name: [{ required: true, message: t('form.nameRequired'), trigger: 'blur' }],
  date: [{ required: true, message: t('form.dateRequired'), trigger: 'change' }],
  message: [{ required: true, message: t('form.messageRequired'), trigger: 'blur' }],
  email: [
    { required: true, message: t('form.emailRequired'), trigger: 'blur' },
    { type: 'email', message: t('form.invalidEmail'), trigger: 'blur' },
  ],
}));

// The backend requires a date strictly after today.
function isNotFutureDate(date) {
  return dayjs(date).endOf('day').isBefore(dayjs().endOf('day').add(1, 'millisecond'));
}

// Detected from the magic bytes: the browser's MIME type comes from the file
// extension and is trivially wrong.
const detectedTypes = new WeakMap();

async function validateAttachment(file) {
  if (file.size > MAX_ATTACHMENT_BYTES) {
    ElMessage.error(t('form.uploadFileSizeError'));
    return false;
  }

  const detected = await fileTypeFromBuffer(await file.arrayBuffer());
  if (!detected || !ALLOWED_IMAGE_TYPES.includes(detected.mime)) {
    ElMessage.error(t('form.uploadFileTypeError'));
    return false;
  }

  detectedTypes.set(file, detected.mime);
  return true;
}

async function uploadAttachmentRequest({ file, onSuccess, onError }) {
  try {
    const contentType = detectedTypes.get(file);
    const target = await getUploadTarget(filesFolderUuid.value, contentType);
    await uploadAttachment(target, file);
    uploadedCount.value += 1;
    onSuccess();
  } catch (error) {
    ElMessage.error(t('form.uploadFailed'));
    onError(error);
  }
}

function handleRemove() {
  uploadedCount.value = Math.max(0, uploadedCount.value - 1);
}

function handleExceed() {
  ElMessage.warning(t('form.uploadLimitExceeded', { count: MAX_ATTACHMENTS }));
}

async function submitForm() {
  const valid = await formRef.value.validate().catch(() => false);
  if (!valid) return;

  isSubmitting.value = true;
  try {
    await createCapsule({
      sender_name: form.name,
      send_at: dayjs(form.date).format('YYYY-MM-DD'),
      message: form.message,
      recipient_email: form.email,
      // Omitted when nothing was uploaded, so the dispatcher skips the fetch.
      ...(uploadedCount.value > 0 ? { files_folder_uuid: filesFolderUuid.value } : {}),
    });
    ElMessage.success(t('form.submitSuccess'));
    resetForm();
  } catch (error) {
    ElMessage.error(error.message);
  } finally {
    isSubmitting.value = false;
  }
}

function resetForm() {
  formRef.value?.resetFields();
  uploadedCount.value = 0;
  filesFolderUuid.value = uuidv4();
}
</script>

<style scoped>
.form-card {
  display: flex;
  max-width: 530px;
  margin: 0 auto;
  margin-top: 50px;
  padding: 10px;
  border-radius: 16px;
  background: radial-gradient(
    circle,
    rgba(41, 123, 134, 0.9),
    rgba(2, 76, 92, 0.8),
    rgba(2, 76, 92, 0.9)
  );
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.5);
  text-align: left;
  flex-direction: column;

  gap: 1.5rem;
  text-align: center;
  border: none;
}

.el-form-item {
  display: flex;
  align-items: center;
  margin-bottom: 20px;
}

.custom-form :deep(.el-form-item__label) {
  color: #dfeefa;
}

.el-form-item .el-form-item__label {
  text-align: right;
  padding-right: 20px;
  width: 150px;
}

.el-input,
.el-date-picker {
  width: 96%;
}

.custom-input {
  width: 96%;
}

.input-field {
  font-size: 14px;
  border-radius: 4px;
  border: 1px solid #ddd;
}

.attachment-container {
  display: flex;
  justify-content: center;
}

::v-deep(.el-upload-list__item-preview) {
  display: none !important;
}

::v-deep(.el-upload-list__item-delete) {
  position: absolute;
  transform: translate(-50%, 0);
}

.form-buttons {
  display: flex;
  justify-content: flex-end;
}

.submit-reset-container {
  justify-content: space-between;
  width: 100%;
}

.submit-button,
.reset-button {
  width: 46%;
  height: 35px;
  font-size: 14px;
  text-align: center;
}

.submit-button {
  background: linear-gradient(45deg, rgba(102, 217, 255, 1) 0%, rgba(45, 99, 255, 1) 100%);
  color: white;
  border: none;
  border-radius: 16px;
  cursor: pointer;
}

.submit-button:hover {
  background: linear-gradient(45deg, rgba(45, 99, 255, 1) 0%, rgba(102, 217, 255, 1) 100%);
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.3);
}

.reset-button {
  background: linear-gradient(45deg, rgba(255, 255, 255, 1) 0%, rgba(234, 234, 234, 1) 100%);
  color: #333;
  border: none;
  border-radius: 16px;
  cursor: pointer;
}

.reset-button:hover {
  background: linear-gradient(45deg, rgba(234, 234, 234, 1) 0%, rgba(255, 255, 255, 1) 100%);
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
}
</style>
