const API_URL = window.location.origin + "/api"

// Login Form Handler
const loginForm = document.getElementById("loginForm")
if (loginForm) {
  loginForm.addEventListener("submit", async (e) => {
    e.preventDefault()
    const email = document.getElementById("email").value
    const password = document.getElementById("password").value

    const response = await fetch(`${API_URL}/auth/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email, password }),
    })

    if (response.ok) {
      const data = await response.json()
      // Store user data in localStorage
      localStorage.setItem("user", JSON.stringify(data.user))
      showNotification("Login berhasil! Mengalihkan...", "success")
      setTimeout(() => {
        window.location.href = "/FrontEnd/dashboard-logged-in.html"
      }, 1500)
    } else {
      showNotification("Email atau password salah", "error")
    }
  })
}

// Register Form Handler
const registerForm = document.getElementById("registerForm")
if (registerForm) {
  registerForm.addEventListener("submit", async (e) => {
    e.preventDefault()
    const email = document.getElementById("email").value
    const password = document.getElementById("password").value
    const confirmPassword = document.getElementById("confirmPassword").value

    // Validate password match
    if (password !== confirmPassword) {
      showNotification("Password tidak cocok!", "error")
      return
    }

    // Validate password length
    if (password.length < 6) {
      showNotification("Password minimal 6 karakter!", "error")
      return
    }

    const response = await fetch(`${API_URL}/auth/register`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        name: email.split("@")[0], // Use email prefix as name
        email,
        password,
        phone: "",
        address: "",
      }),
    })

    if (response.ok) {
      showNotification("Pendaftaran berhasil! Silakan login.", "success")
      setTimeout(() => {
        window.location.href = "/FrontEnd/login.html"
      }, 2000)
    } else {
      showNotification("Pendaftaran gagal. Email mungkin sudah digunakan", "error")
    }
  })
}

const forgetPasswordForm = document.getElementById("forgetPasswordForm")
if (forgetPasswordForm) {
  forgetPasswordForm.addEventListener("submit", async (e) => {
    e.preventDefault()
    const email = document.getElementById("email").value
    const newPassword = document.getElementById("newPassword").value
    const confirmPassword = document.getElementById("confirmPassword").value

    if (!email || email.trim() === "") {
      showNotification("Email harus diisi!", "error")
      return
    }

    if (!newPassword || newPassword.trim() === "") {
      showNotification("Password baru harus diisi!", "error")
      return
    }

    if (!confirmPassword || confirmPassword.trim() === "") {
      showNotification("Konfirmasi password harus diisi!", "error")
      return
    }

    // Validasi password cocok
    if (newPassword !== confirmPassword) {
      showNotification("Password tidak cocok!", "error")
      return
    }

    // Validasi password minimal 6 karakter
    if (newPassword.length < 6) {
      showNotification("Password minimal 6 karakter!", "error")
      return
    }

    const response = await fetch(`${API_URL}/auth/change-password`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email, newPassword }),
    })

    if (response.ok) {
      const data = await response.json()
      if (data.success) {
        showNotification("Password berhasil diubah! Silakan login kembali", "success")
        setTimeout(() => {
          window.location.href = "/FrontEnd/login.html"
        }, 2000)
      } else {
        showNotification("Email tidak ditemukan", "error")
      }
    } else {
      showNotification("Email tidak ditemukan", "error")
    }
  })
}

function showNotification(message, type) {
  const notification = document.createElement("div")
  notification.className = `notification notification-${type}`
  notification.textContent = message
  notification.style.cssText = `
    position: fixed;
    top: 20px;
    right: 20px;
    padding: 15px 20px;
    background-color: ${type === "success" ? "#4CAF50" : "#f44336"};
    color: white;
    border-radius: 4px;
    z-index: 9999;
    animation: slideIn 0.3s ease-in-out;
  `
  document.body.appendChild(notification)

  setTimeout(() => {
    notification.style.animation = "slideOut 0.3s ease-in-out"
    setTimeout(() => notification.remove(), 300)
  }, 3000)
}

// Add animation styles
const style = document.createElement("style")
style.textContent = `
  @keyframes slideIn {
    from {
      transform: translateX(400px);
      opacity: 0;
    }
    to {
      transform: translateX(0);
      opacity: 1;
    }
  }
  @keyframes slideOut {
    from {
      transform: translateX(0);
      opacity: 1;
    }
    to {
      transform: translateX(400px);
      opacity: 0;
    }
  }
`
document.head.appendChild(style)

function logout() {
  localStorage.removeItem("user")
  showNotification("Logout berhasil", "success")
  setTimeout(() => {
    window.location.href = "/FrontEnd/index.html"
  }, 1000)
}

function checkAuth() {
  const user = localStorage.getItem("user")
  if (!user) {
    window.location.href = "/FrontEnd/login.html"
  }
  return JSON.parse(user)
}

function getCurrentUser() {
  const user = localStorage.getItem("user")
  return user ? JSON.parse(user) : null
}