<template>
  <main-layout>
    <el-card class="form-card">

      <!-- Rules Form -->
      <el-form :model="form" :rules="rules" ref="formRef" label-width="150px" class="custom-form">
      
        <!-- Sender's Name -->
      <el-form-item :label="$t('form.senderName')" prop="name">
        <el-input v-model="form.name" maxlength="30" :placeholder="$t('form.namePlaceholder')" />
      </el-form-item>
      
      <!-- Delivery Date -->
      <el-form-item :label="$t('form.deliveryDate')" prop="date">
        <el-date-picker v-model="form.date" type="date" placeholder="Select a date"
        :disabled-date="disablePastDates" />
      </el-form-item>
      
      <!-- Message -->
      <el-form-item :label="$t('form.message')" prop="message">
        <el-input v-model="form.message" type="textarea" maxlength="4096" placeholder="Write your message here" />
      </el-form-item>

      <!-- Recipient Email -->
      <el-form-item :label="$t('form.recipientEmail')" prop="email">
        <el-input v-model="form.email" placeholder="Enter recipient's email" />
      </el-form-item>
      
      <!-- Recipient Telegram -->
      <el-form-item :label="$t('form.recipientTelegram')" prop="telegram">
        <el-input v-model="form.telegram" maxlength="50" placeholder="Enter recipient's Telegram username" />
      </el-form-item>
      
        <!-- Attachments -->
        <!-- action="https://jsonplaceholder.typicode.com/posts/" -->
        <el-form-item :label="$t('form.attachments')">
          <div class="attachment-container">
            <el-upload action="https://jsonplaceholder.typicode.com/posts/" :limit="5" :on-success="handleFileUpload"
              :file-list="form.attachments" accept="image/*" list-type="picture-card" class="file-upload">
              <el-button type="primary" icon="el-icon-upload">
                {{ $t('form.upload') }}
              </el-button>
            </el-upload>
          </div>
        </el-form-item>

        <!-- Submit Button -->
        <el-form-item>
          <el-button type="primary" @click="submitForm">
            {{ $t('form.submit') }}
          </el-button>
          <el-button type="default" @click="resetForm">
            {{ $t('form.reset') }}
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </main-layout>
</template>

<script>
import MainLayout from "@/layouts/MainLayout.vue";
export default {
  components: { MainLayout },
  data() {
    return {
      form: {
        name: "",
        date: "",
        message: "",
        email: "",
        telegram: "",
        attachments: [],
      },
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
    }
  },
  methods: {
    disablePastDates(date) {
      return date.getTime() < Date.now();
    },
    handleFileUpload(response, file, fileList) {
      this.form.attachments = fileList;
    },
    submitForm() {
      this.$refs.formRef.validate((valid) => {
        if (valid) {
          console.log("Form submitted:", this.form);
        }
      });
    },
    resetForm() {
      this.$refs.formRef.resetFields();
      this.form.attachments = [];
    },
  },
};
</script>

<style scoped>
.form-card {
  max-width: 600px;
  margin: 0 auto;
  padding: 20px;
}

.custom-form .el-form-item__label {
  text-align: left;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.attachment-container {
  display: flex;
  align-items: center;
  gap: 10px;
}

.file-upload {
  flex-grow: 1;
}

.el-button {
  margin-right: 10px;
}
</style>