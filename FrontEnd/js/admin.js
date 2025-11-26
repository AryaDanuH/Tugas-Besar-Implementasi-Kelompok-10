// Admin Dashboard JavaScript

const API_BASE = "http://localhost:8080/api"

document.addEventListener("DOMContentLoaded", () => {
  loadDashboardStats()
  loadPendingBorrows()
})

async function loadDashboardStats() {
  try {
    // Load total books
    const booksRes = await fetch(`${API_BASE}/books`)
    const books = await booksRes.json()
    document.getElementById("totalBooks").textContent = books.length || 0

    // Load total users (mock data)
    document.getElementById("totalUsers").textContent = "150"

    // Load pending borrows
    const borrowsRes = await fetch(`${API_BASE}/borrows`)
    const borrows = await borrowsRes.json()
    const pending = borrows.filter((b) => b.status === "pending").length
    const active = borrows.filter((b) => b.status === "approved" || b.status === "borrowed").length

    document.getElementById("pendingBorrows").textContent = pending
    document.getElementById("activeBorrows").textContent = active
  } catch (error) {
    console.error("Error loading stats:", error)
  }
}

async function loadPendingBorrows() {
  try {
    const res = await fetch(`${API_BASE}/borrows`)
    const borrows = await res.json()
    const pending = borrows.filter((b) => b.status === "pending")

    const tbody = document.getElementById("borrowsTable")
    tbody.innerHTML = ""

    if (pending.length === 0) {
      tbody.innerHTML =
        '<tr><td colspan="7" style="text-align: center; color: #999;">Tidak ada peminjaman pending</td></tr>'
      return
    }

    pending.forEach((borrow) => {
      const row = document.createElement("tr")
      row.innerHTML = `
                <td>${borrow.borrow_id}</td>
                <td>User ${borrow.user_id}</td>
                <td>${borrow.book_title}</td>
                <td>${borrow.delivery_type}</td>
                <td>${new Date(borrow.borrow_date).toLocaleDateString("id-ID")}</td>
                <td><span class="status-badge status-pending">${borrow.status}</span></td>
                <td>
                    <div class="action-buttons">
                        <button class="btn-small btn-approve" onclick="approveBorrow(${borrow.borrow_id})">Setujui</button>
                        <button class="btn-small btn-reject" onclick="rejectBorrow(${borrow.borrow_id})">Tolak</button>
                    </div>
                </td>
            `
      tbody.appendChild(row)
    })
  } catch (error) {
    console.error("Error loading borrows:", error)
  }
}

async function approveBorrow(borrowId) {
  try {
    const res = await fetch(`${API_BASE}/borrows/${borrowId}/approve`, {
      method: "PUT",
    })
    const data = await res.json()
    if (data.success) {
      alert("Peminjaman disetujui")
      loadPendingBorrows()
    }
  } catch (error) {
    console.error("Error approving borrow:", error)
    alert("Gagal menyetujui peminjaman")
  }
}

async function rejectBorrow(borrowId) {
  try {
    const res = await fetch(`${API_BASE}/borrows/${borrowId}/reject`, {
      method: "PUT",
    })
    const data = await res.json()
    if (data.success) {
      alert("Peminjaman ditolak")
      loadPendingBorrows()
    }
  } catch (error) {
    console.error("Error rejecting borrow:", error)
    alert("Gagal menolak peminjaman")
  }
}

document.querySelector(".logout-btn").addEventListener("click", () => {
  localStorage.removeItem("user")
  window.location.href = "login.html"
})
