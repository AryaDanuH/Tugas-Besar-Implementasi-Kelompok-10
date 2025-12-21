document.addEventListener("DOMContentLoaded", async () => {
  const checkAuth = () => {
    // Assuming checkAuth function is defined here for the sake of completeness
    // In a real scenario, this function should be imported or defined elsewhere
    const user = JSON.parse(localStorage.getItem("user"))
    return user
  }

  const user = checkAuth()
  if (!user || !user.name) {
    return
  }

  const userNameElement = document.getElementById("userName")
  if (userNameElement) {
    userNameElement.textContent = user.name
  }

  try {
    const borrowResponse = await fetch(`${window.location.origin}/api/borrows/user/${user.user_id}`)
    const borrows = await borrowResponse.json()

    // Calculate stats
    const borrowed = borrows.filter((t) => t.status === "borrowed").length
    const returned = borrows.filter((t) => t.status === "returned").length
    const total = borrows.length

    const borrowedCountElement = document.getElementById("borrowedCount")
    const returnedCountElement = document.getElementById("returnedCount")
    const totalCountElement = document.getElementById("totalCount")

    if (borrowedCountElement) borrowedCountElement.textContent = borrowed
    if (returnedCountElement) returnedCountElement.textContent = returned
    if (totalCountElement) totalCountElement.textContent = total

    // Get locations count
    const locResponse = await fetch(`${window.location.origin}/api/locations`)
    const locations = await locResponse.json()
    const locationCountElement = document.getElementById("locationCount")
    if (locationCountElement) {
      locationCountElement.textContent = locations.length
    }

    // Display recent transactions
    const tbody = document.getElementById("recentTransactions")
    if (tbody) {
      tbody.innerHTML = ""
      borrows.slice(0, 5).forEach((trans) => {
        const row = document.createElement("tr")
        row.innerHTML = `
                <td>${trans.transaction_id}</td>
                <td>${new Date(trans.borrow_date).toLocaleDateString("id-ID")}</td>
                <td>${trans.return_date ? new Date(trans.return_date).toLocaleDateString("id-ID") : "-"}</td>
                <td><span class="status-badge status-${trans.status}">${trans.status}</span></td>
                <td>${trans.delivery_type}</td>
            `
        tbody.appendChild(row)
      })
    }
  } catch (error) {
    console.error("Error loading dashboard:", error)
  }
})
