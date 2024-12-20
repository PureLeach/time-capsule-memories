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
          <el-date-picker v-model="form.date" type="date" :placeholder="$t('form.deliveryDate')"
            :disabled-date="disablePastDates" class="input-field" />
        </el-form-item>

        <!-- Message -->
        <el-form-item :label="$t('form.message')" prop="message">
          <el-input v-model="form.message" type="textarea" maxlength="4096" :placeholder="$t('form.message')"
            class="input-field" />
        </el-form-item>

        <!-- Recipient Email -->
        <el-form-item :label="$t('form.recipientEmail')" prop="email">
          <el-input v-model="form.email" :placeholder="$t('form.recipientEmail')" class="input-field" />
        </el-form-item>

        <!-- Attachments -->
        <el-form-item :label="$t('form.attachments')">
          <div class="attachment-container">
            <el-upload :http-request="uploadToS3" :limit="5" :file-list="form.attachments" accept="image/*"
              list-type="picture-card" class="file-upload" :before-upload="beforeUpload">
              <el-button type="primary" icon="el-icon-upload" class="upload-button">
                {{ $t('form.upload') }}
              </el-button>
            </el-upload>
          </div>
        </el-form-item>

        <!-- Submit Button -->
        <el-form-item>
          <el-button type="primary" @click="submitForm" class="submit-button">
            {{ $t('form.submit') }}
          </el-button>
          <el-button type="default" @click="resetForm" class="reset-button">
            {{ $t('form.reset') }}
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </main-layout>
</template>


<script>
import axios from "axios";
import MainLayout from "@/layouts/MainLayout.vue";
import { v4 as uuidv4 } from "uuid";
import dayjs from "dayjs";
import { fileTypeFromBuffer } from 'file-type'; // Исправленный импорт

const apiUrl = import.meta.env.VITE_BACKEND_API_URL;

axios.defaults.baseURL = apiUrl;

export default {
  components: { MainLayout },
  data() {
    return {
      form: {
        name: "",
        date: "",
        message: "",
        email: "",
      },
      uniqueId: uuidv4(), // Генерация UUID при загрузке страницы
      presignedUrl: "", // Для хранения presigned URL
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
    // Генерация UUID при загрузке страницы
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
        // Отправляем запрос на сервер для получения presigned URL с уникальным uuid
        const response = await axios.get(`/generate-presigned-url?directory=${this.uniqueId}`);
        this.presignedUrl = response.data.presigned_url;
      } catch (error) {
        console.error("Error generating presigned URL:", error);
      }
    },

    // Обработка загрузки файлов на presigned URL
    async uploadToS3({ file }) {
      // Генерация нового presigned URL перед загрузкой файла
      await this.generatePresignedUrl();

      if (!this.presignedUrl) {
        console.error("Presigned URL is missing or invalid");
        return Promise.reject(new Error("Presigned URL is not available"));
      }

      const s3Axios = axios.create();
      try {
        await s3Axios.put(this.presignedUrl, file, {
          headers: {
            "Content-Type": file.type,
          },
        });
        console.log("File uploaded successfully");
      } catch (error) {
        console.error("Error uploading file to S3:", error.message || error);
        throw error;
      }
    },

    // Проверка типа файла перед загрузкой
    async beforeUpload(file) {
      const arrayBuffer = await file.arrayBuffer();
      const type = await fileTypeFromBuffer(arrayBuffer); // Используем fileTypeFromBuffer

      if (!type || !type.mime.startsWith("image/")) {
        this.$message.error("Можно загружать только изображения");
        return false;
      }
      return true;
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
      this.generatePresignedUrl(); // перегенерировать presigned URL для новой формы
    },
  },
};
</script>



<style scoped>
/* Изменить текст лейблов */
.custom-form :deep(.el-form-item__label) {
  color: #dfeefa;
}


/* Фоновая форма */
.form-card {
  max-width: 700px;
  margin: 0 auto;
  padding: 30px;
  background: linear-gradient(135deg, rgba(66, 100, 120, 0.6) 0%, rgba(90, 130, 160, 0.6) 100%);
  border-radius: 25px;
  box-shadow: 0 12px 30px rgba(0, 0, 0, 0.15);
  /* уменьшена непрозрачность тени */
  backdrop-filter: blur(15px);
  border: 2px solid rgba(255, 255, 255, 0.5);
  /* граница с прозрачностью */
  animation: glow 1.5s infinite alternate;

  /* Добавление Flexbox */
  display: flex;
  flex-direction: column;
  align-items: center;
  /* Центрируем элементы по горизонтали */
  justify-content: center;
  /* Центрируем элементы по вертикали */
  gap: 20px;
  /* Добавляем пространство между элементами */
}

@keyframes glow {
  0% {
    box-shadow: 0 0 10px rgba(66, 100, 120, 0.8), 0 0 20px rgba(90, 130, 160, 0.8);
  }

  100% {
    box-shadow: 0 0 20px rgba(66, 100, 120, 0.9), 0 0 40px rgba(90, 130, 160, 0.9);
  }
}


/* Кнопки */

.submit-button {
  background: linear-gradient(45deg, rgba(102, 217, 255, 1) 0%, rgba(45, 99, 255, 1) 100%);
  color: white;
  padding: 10px 20px;
  border: none;
  border-radius: 8px;
  cursor: pointer;
}

.submit-button:hover {
  background: linear-gradient(45deg, rgba(45, 99, 255, 1) 0%, rgba(102, 217, 255, 1) 100%);
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.3);
}

.reset-button {
  background: linear-gradient(45deg, rgba(255, 255, 255, 1) 0%, rgba(234, 234, 234, 1) 100%);
  color: #333;
  padding: 10px 20px;
  border: none;
  border-radius: 8px;
  cursor: pointer;
}

.reset-button:hover {
  background: linear-gradient(45deg, rgba(234, 234, 234, 1) 0%, rgba(255, 255, 255, 1) 100%);
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
}


.el-button {
  margin-right: 10px;
}
</style>


