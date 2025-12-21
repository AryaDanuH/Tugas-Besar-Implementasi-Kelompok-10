let allBooks = []
let selectedBook = null
const API_URL = "http://localhost:8080"

function checkAuth() {
  // Placeholder for authentication check logic
  return { user_id: 1 } // Assuming a user is logged in for demo purposes
}

document.addEventListener("DOMContentLoaded", async () => {
  // Only run on pages that have these containers (not on category-books.html)
  if (document.getElementById("booksContainer") || document.getElementById("categoryFilter")) {
    checkAuth()
    await loadBooks()
    await loadCategories()
  }
})

async function loadBooks() {
  try {
    const response = await fetch(`${API_URL}/api/books`)
    allBooks = await response.json()
    displayBooks(allBooks)
  } catch (error) {
    console.error("Error loading books:", error)
  }
}

async function loadCategories() {
  try {
    const response = await fetch(`${API_URL}/api/books`)
    const books = await response.json()
    const categories = [...new Set(books.map((b) => ({ id: b.category_id, name: b.category_name })))]

    const select = document.getElementById("categoryFilter")
    categories.forEach((cat) => {
      if (cat.id) {
        const option = document.createElement("option")
        option.value = cat.id
        option.textContent = cat.name
        select.appendChild(option)
      }
    })
  } catch (error) {
    console.error("Error loading categories:", error)
  }
}

function displayBooks(books) {
  const container = document.getElementById("booksContainer")
  container.innerHTML = ""

  if (books.length === 0) {
    container.innerHTML = '<p class="text-center">Tidak ada buku ditemukan</p>'
    return
  }

  books.forEach((book) => {
    const card = document.createElement("div")
    card.className = "book-card"
    card.onclick = () => showBookDetails(book)
    card.innerHTML = `
            <div class="book-cover">📚</div>
            <div class="book-info">
                <h3>${book.title}</h3>
                <p><strong>Penulis:</strong> ${book.author}</p>
                <p><strong>Penerbit:</strong> ${book.publisher || "-"}</p>
                <p><strong>Tahun:</strong> ${book.year_published || "-"}</p>
                <span class="book-category">${book.category_name || "Umum"}</span>
            </div>
        `
    container.appendChild(card)
  })
}

function showBookDetails(book) {
  selectedBook = book
  const modal = document.getElementById("bookModal")
  const details = document.getElementById("bookDetails")
  details.innerHTML = `
        <h2>${book.title}</h2>
        <p><strong>Penulis:</strong> ${book.author}</p>
        <p><strong>Penerbit:</strong> ${book.publisher || "-"}</p>
        <p><strong>Tahun Terbit:</strong> ${book.year_published || "-"}</p>
        <p><strong>ISBN:</strong> ${book.isbn || "-"}</p>
        <p><strong>Kategori:</strong> ${book.category_name || "Umum"}</p>
    `
  modal.style.display = "block"
}

function closeModal() {
  document.getElementById("bookModal").style.display = "none"
}

function borrowBook() {
  if (!selectedBook) return
  const user = checkAuth()

  // In a real app, you would select a location and delivery type
  const bookLocationId = 1 // Default for demo
  const deliveryType = "offline"

  fetch(`${API_URL}/api/borrows`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      user_id: user.user_id,
      book_location_id: bookLocationId,
      delivery_type: deliveryType,
    }),
  })
    .then((response) => {
      if (response.ok) {
        alert("Buku berhasil dipinjam!")
        closeModal()
      }
    })
    .catch((error) => console.error("Error:", error))
}

document.getElementById("searchInput")?.addEventListener("input", (e) => {
  const query = e.target.value.toLowerCase()
  const filtered = allBooks.filter(
    (book) => book.title.toLowerCase().includes(query) || book.author.toLowerCase().includes(query),
  )
  displayBooks(filtered)
})

document.getElementById("categoryFilter")?.addEventListener("change", (e) => {
  const categoryId = e.target.value
  if (categoryId) {
    const filtered = allBooks.filter((book) => book.category_id == categoryId)
    displayBooks(filtered)
  } else {
    displayBooks(allBooks)
  }
})
