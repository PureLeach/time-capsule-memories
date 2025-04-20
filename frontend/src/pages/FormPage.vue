<template>
  <main-layout>
    <el-card class="form-card">
      <!-- Rules Form -->
      <el-form :model="form" :rules="rules" ref="formRef" label-width="150px" class="custom-form">
        <!-- Sender's Name -->
        <el-form-item :label="$t('form.senderName')" prop="name">
          <el-input v-model="form.name" maxlength="30" :placeholder="$t('form.namePlaceholder')" class="input-field" />
        </el-form-item>

        <!-- Delivery Date -->
        <el-form-item :label="$t('form.deliveryDate')" prop="date">
          <el-date-picker v-model="form.date" type="date" format="DD/MM/YYYY" :placeholder="$t('form.deliveryDate')"
            style="width: 96%; height: 40px; font-size: 14px;" :disabled-date="disablePastDates" class="input-field" />
        </el-form-item>

        <!-- Message -->
        <el-form-item :label="$t('form.message')" prop="message">
          <el-input v-model="form.message" type="textarea" maxlength="4096" :placeholder="$t('form.message')"
            class="input-field custom-input" />
        </el-form-item>

        <!-- Recipient Email -->
        <el-form-item :label="$t('form.recipientEmail')" prop="email">
          <el-input v-model="form.email" :placeholder="$t('form.recipientEmail')" class="input-field" />
        </el-form-item>

        <!-- Attachments -->
        <el-form-item :label="$t('form.attachments')">
          <div class="attachment-container">
            <el-upload class="file-upload" list-type="picture-card" accept="image/*" :http-request="uploadToS3"
              :limit="3" :before-upload="beforeUpload" :file-list="form.attachments" @exceed="handleExceed">
              <el-icon>
                <Plus />
              </el-icon>
            </el-upload>
          </div>
        </el-form-item>


      </el-form>
      <!-- Submit and Reset Buttons -->
      <el-form-item class="form-buttons">
        <div class="submit-reset-container">
          <el-button type="primary" @click="submitForm" class="submit-button">
            {{ $t('form.submit') }}
          </el-button>
          <el-button type="default" @click="resetForm" class="reset-button">
            {{ $t('form.reset') }}
          </el-button>
        </div>
      </el-form-item>

    </el-card>

  </main-layout>
</template>


<script>
import axios from "axios";
import MainLayout from "@/layouts/MainLayout.vue";
import { v4 as uuidv4 } from "uuid";
import dayjs from "dayjs";
import { fileTypeFromBuffer } from "file-type";
import { Plus } from "@element-plus/icons-vue";

const apiUrl = import.meta.env.VITE_BACKEND_API_URL;

axios.defaults.baseURL = apiUrl;

export default {
  components: { MainLayout, Plus },
  data() {
    return {
      form: {
        name: "",
        date: "",
        message: "",
        email: "",
      },
      uniqueId: uuidv4(),
      presignedUrl: "",
    };
  },
  computed: {
    rules() {
      return {
        name: [{ required: true, message: this.$t("form.nameRequired"), trigger: "blur" }],
        date: [{ required: true, message: this.$t("form.dateRequired"), trigger: "change" }],
        message: [{ required: true, message: this.$t("form.messageRequired"), trigger: "blur" }],
        email: [
          { required: true, message: this.$t("form.emailRequired"), trigger: "blur" },
          { type: "email", message: this.$t("form.invalidEmail"), trigger: "blur" },
        ],
      };
    },
  },
  created() {
    this.uniqueId = uuidv4();
  },
  methods: {
    disablePastDates(date) {
      return date.getTime() < Date.now();
    },
    formatDate(date) {
      return dayjs(date).format("YYYY-MM-DD");
    },
    async generatePresignedUrl() {
      try {
        const response = await axios.get(`/generate-presigned-url?directory=${this.uniqueId}`);
        this.presignedUrl = response.data.presigned_url;
      } catch (error) {
        console.error("Error generating presigned URL:", error);
      }
    },
    async uploadToS3({ file }) {
      await this.generatePresignedUrl();

      if (!this.presignedUrl) {
        console.error("Presigned URL is missing or invalid");
        return Promise.reject(new Error("Presigned URL is not available"));
      }

      const s3Axios = axios.create();
      try {
        const response = await s3Axios.put(this.presignedUrl, file, {
          headers: {
            "Content-Type": file.type,
          },
        });
        console.log("File uploaded successfully", response);
      } catch (error) {
        console.error("Error uploading file to S3:", error.message || error);
        throw error;
      }
    },
    async beforeUpload(file) {
      const arrayBuffer = await file.arrayBuffer();
      const type = await fileTypeFromBuffer(arrayBuffer);

      if (file.size > 5 * 1024 * 1024) {
        this.$message.error(this.$t("form.uploadFileSizeError"));
        return false;
      }

      if (!type || !type.mime.startsWith("image/")) {
        this.$message.error(this.$t("form.uploadFileTypeError"));
        return false;
      }

      return true;
    },
    handleExceed() {
      this.$message.warning(this.$t("form.uploadLimitExceeded"));
    },
    submitForm() {
      this.$refs.formRef.validate((valid) => {
        if (valid) {
          const data = {
            message: this.form.message,
            recipient_email: this.form.email,
            send_at: this.formatDate(this.form.date),
            sender_name: this.form.name,
            files_folder_uuid: this.uniqueId,
          };

          axios
            .post("/capsules", data)
            .then((response) => {
              console.log("Form submitted successfully", response);
              this.resetForm();
            })
            .catch((error) => {
              console.error("Error submitting form:", error);
            });
        }
      });
    },
    resetForm() {
      this.$refs.formRef.resetFields();
      this.form.attachments = [];
      this.generatePresignedUrl();
    },
  },
};
</script>




<style scoped>
/* The main form */
.form-card {
  display: flex;
  max-width: 530px;
  margin: 0 auto;
  margin-top: 50px;
  padding: 10px;
  border-radius: 16px;
  background: radial-gradient(circle, rgba(41, 123, 134, 0.9), rgba(2, 76, 92, 0.8), rgba(2, 76, 92, 0.9));
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.5);
  text-align: left;
  flex-direction: column;

  gap: 1.5rem;
  text-align: center;
  border: none;
}

/* Labels */
.el-form-item {
  display: flex;
  align-items: center;
  margin-bottom: 20px;
}

.custom-form :deep(.el-form-item__label) {
  color: #dfeefa;
}

/* Labels - right alignment */
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


/* Attached files */

.attachment-container {
  display: flex;
  justify-content: center;

}

::v-deep(.el-upload-list__item-preview) {
  display: none !important;
  /* Completely hides the preview icon. */
}

::v-deep(.el-upload-list__item-delete) {
  position: absolute;
  /* Absolute positioning to adjust the location */
  transform: translate(-50%, 0);
  /* Shift for precise centering */
}



/* Buttons */

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
