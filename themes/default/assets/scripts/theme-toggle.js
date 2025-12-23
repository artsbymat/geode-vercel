document.addEventListener("DOMContentLoaded", () => {
  const root = document.documentElement;
  const toggle = document.querySelector(".theme-toggle");
  const syntaxThemeLink = document.getElementById("syntax-theme");
  const STORAGE_KEY = "theme";

  function setTheme(theme) {
    localStorage.setItem("theme", theme);

    document.dispatchEvent(
      new CustomEvent("theme-change", {
        detail: theme,
      }),
    );
  }

  function applyTheme(theme) {
    root.classList.toggle("dark", theme === "dark");
    localStorage.setItem(STORAGE_KEY, theme);

    if (syntaxThemeLink) {
      syntaxThemeLink.href =
        theme === "dark"
          ? "/styles/syntax-dark.css"
          : "/styles/syntax-light.css";
    }
  }

  document.addEventListener("theme-change", (event) => {
    applyTheme(event.detail);
  });

  toggle.addEventListener("click", () => {
    const isDark = root.classList.contains("dark");
    setTheme(isDark ? "light" : "dark");
  });
});
