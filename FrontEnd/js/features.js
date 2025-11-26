function showFeature(featureName) {
  // Hide all features
  document.querySelectorAll("[data-feature]").forEach((el) => {
    el.style.display = "none"
  })

  // Show selected feature
  const feature = document.querySelector(`[data-feature="${featureName}"]`)
  if (feature) {
    feature.style.display = "block"
    window.scrollTo(0, 0)
  }

  // Update active navigation link
  document.querySelectorAll("[data-nav-link]").forEach((link) => {
    link.classList.remove("active")
  })
  const activeLink = document.querySelector(`[data-nav-link="${featureName}"]`)
  if (activeLink) {
    activeLink.classList.add("active")
  }
}

function initializeFeatures() {
  // Set default view to home
  showFeature("home")

  // Add click handlers to navigation links
  document.querySelectorAll("[data-nav-link]").forEach((link) => {
    link.addEventListener("click", (e) => {
      e.preventDefault()
      const feature = link.getAttribute("data-nav-link")
      showFeature(feature)
    })
  })
}

document.addEventListener("DOMContentLoaded", initializeFeatures)
