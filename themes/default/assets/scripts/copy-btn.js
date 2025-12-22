(function () {
  function getCodeText(codeEl) {
    return codeEl.innerText.replace(/\n$/, "");
  }

  function copyCode(button) {
    const block = button.closest(".code-block");
    if (!block) return;

    const code = block.querySelector("pre > code");
    if (!code) return;

    const text = getCodeText(code);

    navigator.clipboard
      .writeText(text)
      .then(() => {
        const original = button.textContent;
        button.textContent = "Copied!";
        button.disabled = true;

        setTimeout(() => {
          button.textContent = original;
          button.disabled = false;
        }, 1500);
      })
      .catch(() => {
        const textarea = document.createElement("textarea");
        textarea.value = text;
        textarea.style.position = "fixed";
        textarea.style.opacity = "0";
        document.body.appendChild(textarea);
        textarea.select();
        document.execCommand("copy");
        document.body.removeChild(textarea);
      });
  }

  document.addEventListener("click", function (e) {
    const btn = e.target.closest(".copy-btn");
    if (!btn) return;
    copyCode(btn);
  });
})();
