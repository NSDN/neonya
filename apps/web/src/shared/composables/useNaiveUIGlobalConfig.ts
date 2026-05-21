import {
  darkTheme,
  lightTheme,
  useMessage,
  type GlobalTheme
} from "naive-ui";
import { ref, computed } from "vue";

export function useNaiveUIGlobalConfig() {
  const currentTheme = ref<string>("dark");

  const theme = computed<GlobalTheme>(() =>
    currentTheme.value === "dark" ? darkTheme : lightTheme
  );

  const toggleTheme = () => {
    currentTheme.value = currentTheme.value === "dark" ? "light" : "dark";
  };

  const initMessager = () => {
    window.$message = useMessage();
  };

  return {
    theme,
    toggleTheme,
    initMessager,
  };
}
