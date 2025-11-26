// Import the checkAuth function
function checkAuth() {
  // Placeholder for actual authentication check logic
  return { name: "John Doe", user_id: 123 }
}

const API_URL = window.location.origin + "/api"

document.addEventListener("DOMContentLoaded", async () => {
  const user = checkAuth()
  document.getElementById("userName").textContent = user.name

  try {
    const borrowResponse = await fetch(`${API_URL}/borrows/user/${user.user_id}`)
    const borrows = await borrowResponse.json()

    // Calculate stats
    const borrowed = borrows.filter((t) => t.status === "borrowed").length
    const returned = borrows.filter((t) => t.status === "returned").length
    const total = borrows.length

    document.getElementById("borrowedCount").textContent = borrowed
    document.getElementById("returnedCount").textContent = returned
    document.getElementById("totalCount").textContent = total

    // Get locations count
    const locResponse = await fetch(`${API_URL}/locations`)
    const locations = await locResponse.json()
    document.getElementById("locationCount").textContent = locations.length

    // Display recent transactions
    const tbody = document.getElementById("recentTransactions")
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
  } catch (error) {
    console.error("Error loading dashboard:", error)
  }
})
