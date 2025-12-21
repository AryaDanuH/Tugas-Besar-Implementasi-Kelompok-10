// Borrow Management JavaScript

const API_BASE = "http://localhost:8080/api"
let currentTab = "pending"

document.addEventListener("DOMContentLoaded", () => {
  loadBorrowsByStatus("pending")
})

function switchTab(tab) {
  currentTab = tab
  document.querySelectorAll(".tab-btn").forEach((btn) => btn.classList.remove("active"))
  document.querySelectorAll(".tab-content").forEach((content) => content.classList.remove("active"))

  event.target.classList.add("active")
  document.getElementById(tab).classList.add("active")

  loadBorrowsByStatus(tab)
}

async function loadBorrowsByStatus(status) {
  try {
    const res = await fetch(`${API_BASE}/borrows`)
    const borrows = await res.json()
    const filtered = borrows.filter((b) => b.status === status)

    const tableId = `${status}Table`
    const tbody = document.getElementById(tableId)
    tbody.innerHTML = ""

    if (filtered.length === 0) {
      const cols = status === "pending" || status === "approved" ? 6 : status === "borrowed" ? 5 : 5
      tbody.innerHTML = `<tr><td colspan="${cols}" style="text-align: center; color: #999;">Tidak ada data</td></tr>`
      return
    }

    filtered.forEach((borrow) => {
      const row = document.createElement("tr")
      let html = `
                <td>${borrow.borrow_id}</td>
                <td>User ${borrow.user_id}</td>
                <td>${borrow.book_title}</td>
            `

      if (status === "pending" || status === "approved") {
        html += `
                    <td>${borrow.delivery_type}</td>
                    <td>${new Date(borrow.borrow_date).toLocaleDateString("id-ID")}</td>
                    <td>
                        <div class="action-buttons">
                            ${
                              status === "pending"
                                ? `
                                <button class="btn-small btn-approve" onclick="approveBorrow(${borrow.borrow_id})">Setujui</button>
                                <button class="btn-small btn-reject" onclick="rejectBorrow(${borrow.borrow_id})">Tolak</button>
                            `
                                : `
                                <button class="btn-small btn-return" onclick="returnBook(${borrow.borrow_id})">Kembalikan</button>
                            `
                            }
                        </div>
                    </td>
                `
      } else if (status === "borrowed") {
        html += `
                    <td>${new Date(borrow.borrow_date).toLocaleDateString("id-ID")}</td>
                    <td>
                        <button class="btn-small btn-return" onclick="returnBook(${borrow.borrow_id})">Kembalikan</button>
                    </td>
                `
      } else {
        html += `
                    <td>${new Date(borrow.borrow_date).toLocaleDateString("id-ID")}</td>
                    ${borrow.return_date ? `<td>${new Date(borrow.return_date).toLocaleDateString("id-ID")}</td>` : "<td>-</td>"}
                `
      }

      row.innerHTML = html
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
      loadBorrowsByStatus(currentTab)
    }
  } catch (error) {
    console.error("Error:", error)
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
      loadBorrowsByStatus(currentTab)
    }
  } catch (error) {
    console.error("Error:", error)
    alert("Gagal menolak peminjaman")
  }
}

async function returnBook(borrowId) {
  try {
    const res = await fetch(`${API_BASE}/borrows/${borrowId}/return`, {
      method: "PUT",
    })
    const data = await res.json()
    if (data.success) {
      alert("Buku berhasil dikembalikan")
      loadBorrowsByStatus(currentTab)
    }
  } catch (error) {
    console.error("Error:", error)
    alert("Gagal mengembalikan buku")
  }
}

document.querySelector(".logout-btn").addEventListener("click", () => {
  localStorage.removeItem("user")
  window.location.href = "login.html"
})
