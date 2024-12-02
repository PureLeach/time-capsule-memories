<template>
  <el-select
    v-model="currentLocale"
    size="small"
    @change="changeLanguage"
    class="language-switcher"
    placeholder="Language"
  >
    <el-option label="English" value="en" />
    <el-option label="Русский" value="ru" />
  </el-select>
</template>

<script>
import { useI18n } from 'vue-i18n';

export default {
  setup() {
    const { locale } = useI18n();
    const currentLocale = locale;

    const changeLanguage = (value) => {
      locale.value = value;
    };

    return { currentLocale, changeLanguage };
  },
};
</script>

<style scoped>
.language-switcher {
  position: fixed;
  top: 10px;
  right: 10px;
  background-color: #fff;
  border-radius: 5px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  z-index: 9999;
  width: auto;  /* Убираем фулл ширину */
  min-width: 120px;  /* Минимальная ширина */
}

.language-switcher .el-select {
  width: 100%;
  max-width: 120px;  /* Максимальная ширина */
  font-size: 14px;
}

.language-switcher .el-option {
  font-size: 14px;
}

.language-switcher .el-select-dropdown {
  border-radius: 5px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.language-switcher .el-select .el-input__inner {
  font-size: 14px;
  padding: 5px 10px;
}

.language-switcher .el-select .el-input__inner:hover {
  border-color: #409eff;
}
</style>
