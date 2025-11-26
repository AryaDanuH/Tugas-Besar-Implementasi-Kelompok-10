// Book Management JavaScript

const API_BASE = "http://localhost:8080/api"
let editingBookId = null

document.addEventListener("DOMContentLoaded", () => {
  loadBooks()
})

async function loadBooks() {
  try {
    const res = await fetch(`${API_BASE}/books`)
    const books = await res.json()
    displayBooks(books)
  } catch (error) {
    console.error("Error loading books:", error)
    document.getElementById("booksGrid").innerHTML =
      '<p style="grid-column: 1/-1; text-align: center; color: #999;">Gagal memuat buku</p>'
  }
}

function displayBooks(books) {
  const grid = document.getElementById("booksGrid")
  grid.innerHTML = ""

  if (books.length === 0) {
    grid.innerHTML = '<p style="grid-column: 1/-1; text-align: center; color: #999;">Tidak ada buku</p>'
    return
  }

  books.forEach((book) => {
    const card = document.createElement("div")
    card.className = "book-item"
    card.innerHTML = `
            <div class="book-item-image">Sampul Buku</div>
            <div class="book-item-content">
                <div class="book-item-title">${book.title}</div>
                <div class="book-item-author">${book.author}</div>
                <div class="book-item-actions">
                    <button class="btn-small btn-edit" onclick="editBook(${book.book_id})">Edit</button>
                    <button class="btn-small btn-delete" onclick="deleteBook(${book.book_id})">Hapus</button>
                </div>
            </div>
        `
    grid.appendChild(card)
  })
}

function openAddBookModal() {
  editingBookId = null
  document.getElementById("bookForm").reset()
  document.getElementById("bookModal").classList.add("active")
}

function closeBookModal() {
  document.getElementById("bookModal").classList.remove("active")
}

function editBook(bookId) {
  editingBookId = bookId
  // Load book data and populate form
  document.getElementById("bookModal").classList.add("active")
}

async function deleteBook(bookId) {
  if (!confirm("Yakin ingin menghapus buku ini?")) return

  try {
    const res = await fetch(`${API_BASE}/books/${bookId}`, {
      method: "DELETE",
    })
    const data = await res.json()
    if (data.success) {
      alert("Buku berhasil dihapus")
      loadBooks()
    }
  } catch (error) {
    console.error("Error deleting book:", error)
    alert("Gagal menghapus buku")
  }
}

function filterBooks() {
  const search = document.getElementById("searchInput").value
  const category = document.getElementById("categoryFilter").value
  // Implement filter logic
  loadBooks()
}

document.getElementById("bookForm").addEventListener("submit", async (e) => {
  e.preventDefault()

  const bookData = {
    title: document.getElementById("bookTitle").value,
    author: document.getElementById("bookAuthor").value,
    publisher: document.getElementById("bookPublisher").value,
    year_published: Number.parseInt(document.getElementById("bookYear").value),
    isbn: document.getElementById("bookISBN").value,
    category_id: Number.parseInt(document.getElementById("bookCategory").value),
  }

  try {
    let res
    if (editingBookId) {
      res = await fetch(`${API_BASE}/books/${editingBookId}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(bookData),
      })
    } else {
      res = await fetch(`${API_BASE}/books`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(bookData),
      })
    }

    const data = await res.json()
    if (data.success) {
      alert(editingBookId ? "Buku berhasil diperbarui" : "Buku berhasil ditambahkan")
      closeBookModal()
      loadBooks()
    }
  } catch (error) {
    console.error("Error saving book:", error)
    alert("Gagal menyimpan buku")
  }
})

document.querySelector(".logout-btn").addEventListener("click", () => {
  localStorage.removeItem("user")
  window.location.href = "login.html"
})
