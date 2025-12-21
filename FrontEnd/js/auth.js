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
      showNotification("Login successful! Redirecting...", "success")
      setTimeout(() => {
        if (data.user.role === "admin") {
          window.location.href = "/FrontEnd/admin-all-books.html"
        } else {
          window.location.href = "/FrontEnd/dashboard-logged-in.html"
        }
      }, 1500)
    } else {
      showNotification("Incorrect email or password", "error")
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
      showNotification("Passwords do not match!", "error")
      return
    }

    // Validate password length
    if (password.length < 6) {
      showNotification("Password must be at least 6 characters!", "error")
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
      showNotification("Registration successful! Please login.", "success")
      setTimeout(() => {
        window.location.href = "/FrontEnd/login.html"
      }, 2000)
    } else {
      showNotification("Registration failed. Email may already be in use", "error")
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
      showNotification("Email is required!", "error")
      return
    }

    if (!newPassword || newPassword.trim() === "") {
      showNotification("New password is required!", "error")
      return
    }

    if (!confirmPassword || confirmPassword.trim() === "") {
      showNotification("Password confirmation is required!", "error")
      return
    }

    // Validasi password cocok
    if (newPassword !== confirmPassword) {
      showNotification("Passwords do not match!", "error")
      return
    }

    // Validasi password minimal 6 karakter
    if (newPassword.length < 6) {
      showNotification("Password must be at least 6 characters!", "error")
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
        showNotification("Password successfully changed! Please login again", "success")
        setTimeout(() => {
          window.location.href = "/FrontEnd/login.html"
        }, 2000)
      } else {
        showNotification("Email not found", "error")
      }
    } else {
      showNotification("Email not found", "error")
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
  showNotification("Logout successful", "success")
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

function loadUserProfileImage() {
  const user = getCurrentUser()
  if (!user) return

  // Find all profile avatar images on the page
  const avatarImages = document.querySelectorAll("#userAvatar, #userInitial, .profile-avatar-img")

  avatarImages.forEach((img) => {
    // If user has a profile image, use it
    if (user.profile_image && user.profile_image.trim() !== "") {
      img.src = user.profile_image
      img.alt = user.name || "User Profile"
    } else {
      // Use placeholder if no profile image
      img.src = "/FrontEnd/images/placeholder-profile.png"
      img.alt = user.name ? user.name.charAt(0).toUpperCase() : "U"
    }
  })

  // Update profile image preview on profile page
  const profilePreview = document.getElementById("profileImagePreview")
  if (profilePreview) {
    if (user.profile_image && user.profile_image.trim() !== "") {
      profilePreview.src = user.profile_image
    } else {
      profilePreview.src = "/FrontEnd/images/placeholder-profile.png"
    }
  }
}

// Call loadUserProfileImage when page loads
if (document.readyState === "loading") {
  document.addEventListener("DOMContentLoaded", loadUserProfileImage)
} else {
  loadUserProfileImage()
}
