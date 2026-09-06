<template>
  <main-layout>
    <warp-overlay ref="warp" />

    <section class="console">
      <header class="console-head">
        <span class="eyebrow">{{ t('form.eyebrow') }}</span>
        <h1 class="console-title">{{ t('form.title') }}</h1>
      </header>

      <div class="console-body">
        <div class="panel console-form">
          <el-form
            ref="formRef"
            :model="form"
            :rules="rules"
            label-position="top"
            class="capsule-form"
          >
            <div class="field-row">
              <el-form-item :label="labelFor('01', t('form.senderName'))" prop="name">
                <el-input
                  v-model="form.name"
                  maxlength="100"
                  :placeholder="t('form.namePlaceholder')"
                />
              </el-form-item>

              <el-form-item :label="labelFor('02', t('form.recipientEmail'))" prop="email">
                <el-input
                  v-model="form.email"
                  maxlength="255"
                  :placeholder="t('form.emailPlaceholder')"
                />
              </el-form-item>
            </div>

            <el-form-item :label="labelFor('03', t('form.deliveryDate'))" prop="date">
              <el-date-picker
                v-model="form.date"
                type="date"
                format="DD MMM YYYY"
                :placeholder="t('form.pickDate')"
                :disabled-date="isNotFutureDate"
                :shortcuts="shortcuts"
              />
            </el-form-item>

            <el-form-item :label="labelFor('04', t('form.message'))" prop="message">
              <el-input
                v-model="form.message"
                type="textarea"
                :rows="6"
                maxlength="4096"
                show-word-limit
                resize="none"
                :placeholder="t('form.messagePlaceholder')"
              />
            </el-form-item>

            <el-form-item :label="labelFor('05', t('form.attachments'))">
              <el-upload
                class="uploader"
                list-type="picture-card"
                accept="image/*"
                :limit="MAX_ATTACHMENTS"
                :http-request="uploadAttachmentRequest"
                :before-upload="validateAttachment"
                :on-remove="handleRemove"
                :on-exceed="handleExceed"
              >
                <span class="uploader-plus">+</span>
              </el-upload>
              <p class="mono uploader-hint">{{ t('form.attachmentsHint') }}</p>
            </el-form-item>

            <div class="actions">
              <button type="button" class="btn btn-ghost" @click="resetForm">
                {{ t('form.reset') }}
              </button>
              <button
                type="button"
                class="btn btn-primary"
                :disabled="isSubmitting"
                @click="submitForm"
              >
                <span v-if="isSubmitting" class="spinner"></span>
                {{ isSubmitting ? t('form.sealing') : t('form.submit') }}
              </button>
            </div>
          </el-form>
        </div>

        <aside class="console-side">
          <div class="panel side-card">
            <capsule-preview
              :sender="form.name"
              :recipient="form.email"
              :message="form.message"
              :attachments="uploadedCount"
              :sealed="isSubmitting"
            />
          </div>

          <div class="panel side-card">
            <countdown-readout :target="form.date" />
          </div>
        </aside>
      </div>
    </section>

    <transition name="fade">
      <div v-if="sealed" class="dialog-mask" @click.self="closeSealed">
        <div class="panel dialog">
          <div class="dialog-orb"></div>
          <span class="eyebrow">{{ t('form.sealedTag') }}</span>
          <h2 class="dialog-title">{{ t('form.sealedTitle') }}</h2>
          <p class="dialog-text">
            {{ t('form.sealedText', { date: sealed.date, email: sealed.email }) }}
          </p>
          <dl class="dialog-meta">
            <div>
              <dt class="mono">{{ t('form.sealedId') }}</dt>
              <dd class="mono">#{{ String(sealed.id).padStart(6, '0') }}</dd>
            </div>
            <div>
              <dt class="mono">{{ t('form.sealedStatus') }}</dt>
              <dd class="mono accent">{{ sealed.status }}</dd>
            </div>
          </dl>
          <button type="button" class="btn btn-primary" @click="closeSealed">
            {{ t('form.sealedClose') }}
          </button>
        </div>
      </div>
    </transition>
  </main-layout>
</template>

<script setup>
import { computed, reactive, ref } from 'vue';
import { v4 as uuidv4 } from 'uuid';
import dayjs from 'dayjs';
import { fileTypeFromBuffer } from 'file-type';
import { ElMessage } from 'element-plus';
import { useI18n } from 'vue-i18n';
import MainLayout from '@/layouts/MainLayout.vue';
import WarpOverlay from '@/components/WarpOverlay.vue';
import CapsulePreview from '@/components/CapsulePreview.vue';
import CountdownReadout from '@/components/CountdownReadout.vue';
import { createCapsule, getUploadTarget, uploadAttachment } from '@/api/capsules';

// Mirrors the signed upload policy. Checking here only buys faster feedback.
const MAX_ATTACHMENTS = 3;
const MAX_ATTACHMENT_BYTES = 5 * 1024 * 1024;
const ALLOWED_IMAGE_TYPES = ['image/jpeg', 'image/png', 'image/webp', 'image/gif'];

const { t } = useI18n();

const formRef = ref(null);
const warp = ref(null);
const isSubmitting = ref(false);
const uploadedCount = ref(0);
const sealed = ref(null);
// One folder per capsule links its uploads to the row, without the backend ever
// seeing the files.
const filesFolderUuid = ref(uuidv4());

const form = reactive({ name: '', date: '', message: '', email: '' });

const labelFor = (index, text) => `${index} · ${text}`;

const rules = computed(() => ({
  name: [{ required: true, message: t('form.nameRequired'), trigger: 'blur' }],
  date: [{ required: true, message: t('form.dateRequired'), trigger: 'change' }],
  message: [{ required: true, message: t('form.messageRequired'), trigger: 'blur' }],
  email: [
    { required: true, message: t('form.emailRequired'), trigger: 'blur' },
    { type: 'email', message: t('form.invalidEmail'), trigger: 'blur' },
  ],
}));

const shortcuts = computed(() =>
  [
    { key: 'year', add: [1, 'year'] },
    { key: 'fiveYears', add: [5, 'year'] },
    { key: 'decade', add: [10, 'year'] },
  ].map(({ key, add }) => ({
    text: t(`form.shortcuts.${key}`),
    value: () =>
      dayjs()
        .add(...add)
        .toDate(),
  }))
);

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
    const target = await getUploadTarget(filesFolderUuid.value, detectedTypes.get(file));
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
    const { data } = await createCapsule({
      sender_name: form.name,
      send_at: dayjs(form.date).format('YYYY-MM-DD'),
      message: form.message,
      recipient_email: form.email,
      // Omitted when nothing was uploaded, so the dispatcher skips the fetch.
      ...(uploadedCount.value > 0 ? { files_folder_uuid: filesFolderUuid.value } : {}),
    });

    warp.value?.trigger();
    sealed.value = {
      id: data.id,
      status: data.status,
      email: data.recipient_email,
      date: dayjs(data.send_at).format('DD MMM YYYY'),
    };
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

function closeSealed() {
  sealed.value = null;
}
</script>

<style scoped>
.console {
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.console-head {
  display: flex;
  flex-direction: column;
  gap: 0.8rem;
}

.console-title {
  font-size: clamp(1.8rem, 3.4vw, 2.6rem);
  font-weight: 300;
  letter-spacing: -0.01em;
}

.console-body {
  display: grid;
  grid-template-columns: minmax(0, 1.6fr) minmax(260px, 0.85fr);
  gap: 1.5rem;
  align-items: start;
}

.console-form {
  padding: clamp(1.4rem, 3vw, 2.2rem);
}

.field-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1.2rem;
}

.capsule-form :deep(.el-form-item__label) {
  font-family: var(--font-mono);
  font-size: 0.62rem;
  letter-spacing: 0.18em;
  text-transform: uppercase;
  color: var(--ink-faint);
  padding-bottom: 0.45rem;
}

.capsule-form :deep(.el-date-editor.el-input) {
  width: 100%;
}

.capsule-form :deep(.el-form-item__content) {
  display: block;
}

.uploader-hint {
  margin-top: 0.5rem;
  font-size: 0.58rem;
  text-transform: none;
  letter-spacing: 0.06em;
}

.uploader :deep(.el-upload--picture-card),
.uploader :deep(.el-upload-list__item) {
  width: 96px;
  height: 96px;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px dashed var(--line);
  transition:
    border-color 0.25s ease,
    box-shadow 0.25s ease;
}

.uploader :deep(.el-upload--picture-card:hover) {
  border-color: var(--aqua);
  box-shadow: 0 0 22px rgba(94, 242, 224, 0.22);
}

.uploader :deep(.el-upload-list__item) {
  border-style: solid;
  overflow: hidden;
}

.uploader :deep(.el-upload-list__item-preview) {
  display: none;
}

.uploader-plus {
  font-size: 1.6rem;
  font-weight: 200;
  color: var(--aqua);
}

.actions {
  display: flex;
  justify-content: flex-end;
  gap: 0.8rem;
  margin-top: 0.4rem;
}

.btn {
  display: inline-flex;
  align-items: center;
  gap: 0.6rem;
  padding: 0.75rem 1.6rem;
  border-radius: 999px;
  font-family: var(--font-mono);
  font-size: 0.68rem;
  letter-spacing: 0.18em;
  text-transform: uppercase;
  cursor: pointer;
  transition: all 0.25s ease;
}

.btn-ghost {
  border: 1px solid var(--line);
  background: transparent;
  color: var(--ink-soft);
}

.btn-ghost:hover {
  border-color: var(--ink-faint);
  color: var(--ink);
}

.btn-primary {
  border: 1px solid var(--aqua);
  background: rgba(94, 242, 224, 0.12);
  color: var(--ink);
}

.btn-primary:hover:not(:disabled) {
  background: rgba(94, 242, 224, 0.2);
  box-shadow: 0 0 26px rgba(94, 242, 224, 0.3);
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: progress;
}

.spinner {
  width: 12px;
  height: 12px;
  border: 1.5px solid var(--aqua);
  border-top-color: transparent;
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
}

.console-side {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
  position: sticky;
  top: 1.5rem;
}

.side-card {
  padding: 1.4rem;
}

.dialog-mask {
  position: fixed;
  inset: 0;
  z-index: 2200;
  display: grid;
  place-items: center;
  padding: 1.5rem;
  background: rgba(4, 5, 12, 0.8);
  backdrop-filter: blur(8px);
}

.dialog {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
  width: min(430px, 100%);
  padding: 2.4rem 2rem;
  text-align: center;
  overflow: hidden;
}

.dialog-orb {
  position: absolute;
  top: -70px;
  left: 50%;
  width: 190px;
  height: 190px;
  transform: translateX(-50%);
  border-radius: 50%;
  background: radial-gradient(circle, rgba(255, 196, 107, 0.3), transparent 68%);
  pointer-events: none;
}

.dialog-title {
  font-size: 1.35rem;
  font-weight: 300;
}

.dialog-text {
  font-size: 0.88rem;
  line-height: 1.6;
  color: var(--ink-soft);
}

.dialog-meta {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1px;
  width: 100%;
  background: var(--line);
  border: 1px solid var(--line);
  border-radius: 10px;
  overflow: hidden;
}

.dialog-meta > div {
  padding: 0.7rem;
  background: rgba(8, 10, 22, 0.8);
}

.dialog-meta dt {
  font-size: 0.52rem;
  margin-bottom: 0.25rem;
}

.dialog-meta dd {
  font-size: 0.78rem;
  color: var(--ink);
  text-transform: none;
}

.dialog-meta .accent {
  color: var(--amber);
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.35s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

@media (max-width: 980px) {
  .console-body {
    grid-template-columns: 1fr;
  }
  .console-side {
    position: static;
    flex-direction: row;
    flex-wrap: wrap;
  }
  .side-card {
    flex: 1 1 260px;
  }
}

@media (max-width: 620px) {
  .field-row {
    grid-template-columns: 1fr;
  }
  .actions {
    flex-direction: column-reverse;
  }
  .btn {
    justify-content: center;
  }
}
</style>
