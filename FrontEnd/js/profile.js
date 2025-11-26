const API_URL = window.location.origin + "/api"

// Function to get current user from localStorage
function getCurrentUser() {
  const userString = localStorage.getItem("user")
  return userString ? JSON.parse(userString) : null
}

// Function to show notification
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

document.addEventListener("DOMContentLoaded", () => {
  const user = getCurrentUser()
  if (!user) {
    window.location.href = "/FrontEnd/login.html"
    return
  }

  // Display user data from localStorage
  document.getElementById("name").value = user.name || ""
  document.getElementById("email").value = user.email || ""
  document.getElementById("phone").value = user.phone || ""
  document.getElementById("address").value = user.address || ""
  document.getElementById("role").value = user.role || "user"
  document.getElementById("userName").textContent = user.name || user.email

  if (user.profile_image) {
    document.getElementById("profileImagePreview").src = user.profile_image
  }

  const profileImageInput = document.getElementById("profileImageInput")
  profileImageInput.addEventListener("change", async (e) => {
    const file = e.target.files[0]
    if (!file) return

    // Show preview
    const reader = new FileReader()
    reader.onload = (event) => {
      document.getElementById("profileImagePreview").src = event.target.result
    }
    reader.readAsDataURL(file)

    // Upload file
    const formData = new FormData()
    formData.append("file", file)

    try {
      const response = await fetch(`${API_URL}/users/${user.user_id}/upload-profile-image`, {
        method: "POST",
        body: formData,
      })

      const data = await response.json()

      if (response.ok) {
        // Update user data in localStorage
        const updatedUser = { ...user, profile_image: data.profile_image }
        localStorage.setItem("user", JSON.stringify(updatedUser))
        showNotification("Foto profil berhasil diperbarui!", "success")
      } else {
        showNotification("Gagal mengunggah foto profil", "error")
      }
    } catch (error) {
      console.error("Error uploading profile image:", error)
      showNotification("Terjadi kesalahan: " + error.message, "error")
    }
  })
})

document.getElementById("profileForm").addEventListener("submit", async (e) => {
  e.preventDefault()
  const user = getCurrentUser()

  const updatedUser = {
    name: document.getElementById("name").value,
    phone: document.getElementById("phone").value,
    address: document.getElementById("address").value,
  }

  try {
    const response = await fetch(`${API_URL}/users/${user.user_id}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(updatedUser),
    })

    if (response.ok) {
      const updatedUserData = { ...user, ...updatedUser }
      localStorage.setItem("user", JSON.stringify(updatedUserData))
      showNotification("Profil berhasil diperbarui!", "success")
    } else {
      showNotification("Gagal memperbarui profil", "error")
    }
  } catch (error) {
    console.error("Error updating profile:", error)
    showNotification("Terjadi kesalahan: " + error.message, "error")
  }
})
