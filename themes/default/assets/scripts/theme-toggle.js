document.addEventListener("DOMContentLoaded", () => {
  const root = document.documentElement;
  const toggle = document.querySelector(".theme-toggle");
  const STORAGE_KEY = "theme";

  function getPreferredTheme() {
    const stored = localStorage.getItem(STORAGE_KEY);
    if (stored === "light" || stored === "dark") {
      return stored;
    }

    return window.matchMedia("(prefers-color-scheme: dark)").matches
      ? "dark"
      : "light";
  }

  function applyTheme(theme) {
    root.classList.toggle("dark", theme === "dark");
    localStorage.setItem(STORAGE_KEY, theme);
  }

  applyTheme(getPreferredTheme());

  toggle.addEventListener("click", () => {
    const isDark = root.classList.contains("dark");
    applyTheme(isDark ? "light" : "dark");
  });
});
