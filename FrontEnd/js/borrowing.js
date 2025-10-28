const API_URL = window.location.origin + "/api"
let allTransactions = []

function checkAuth() {
  // Placeholder for authentication check logic
  return { user_id: 1 } // Example user ID
}

document.addEventListener("DOMContentLoaded", async () => {
  const user = checkAuth()
  await loadTransactions(user.user_id)
})

async function loadTransactions(userId) {
  try {
    const response = await fetch(`${API_URL}/borrows/user/${userId}`)
    allTransactions = await response.json()
    displayTransactions(allTransactions)
  } catch (error) {
    console.error("Error loading transactions:", error)
  }
}

function displayTransactions(transactions) {
  const tbody = document.getElementById("borrowingsTable")
  tbody.innerHTML = ""

  if (transactions.length === 0) {
    tbody.innerHTML = '<tr><td colspan="6" class="text-center">Tidak ada peminjaman</td></tr>'
    return
  }

  transactions.forEach((trans) => {
    const row = document.createElement("tr")
    const statusClass = `status-${trans.status}`
    row.innerHTML = `
            <td>${trans.transaction_id}</td>
            <td>${new Date(trans.borrow_date).toLocaleDateString("id-ID")}</td>
            <td>${trans.return_date ? new Date(trans.return_date).toLocaleDateString("id-ID") : "-"}</td>
            <td><span class="status-badge ${statusClass}">${trans.status}</span></td>
            <td>${trans.delivery_type}</td>
            <td>
                ${trans.status === "borrowed" ? `<button class="btn btn-success" onclick="returnTransaction(${trans.transaction_id})">Kembalikan</button>` : "-"}
            </td>
        `
    tbody.appendChild(row)
  })
}

function filterTransactions(status) {
  const buttons = document.querySelectorAll(".tab-btn")
  buttons.forEach((btn) => btn.classList.remove("active"))
  event.target.classList.add("active")

  if (status === "all") {
    displayTransactions(allTransactions)
  } else {
    const filtered = allTransactions.filter((t) => t.status === status)
    displayTransactions(filtered)
  }
}

function returnTransaction(transactionId) {
  fetch(`${API_URL}/borrows/${transactionId}/return`, {
    method: "PUT",
  })
    .then((response) => {
      if (response.ok) {
        alert("Buku berhasil dikembalikan!")
        const user = checkAuth()
        loadTransactions(user.user_id)
      }
    })
    .catch((error) => console.error("Error:", error))
}
