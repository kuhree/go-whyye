const selectors = {
  form: "#form--name",
  quote: "#quote"
}

async function handleFormSubmit(e) {
  e.stopPropagation()
  e.preventDefault()

  const { target } = e;

  if (target && !target.classList.contains("loading")) {
    target.classList.add("loading")
  }

  const formData = new FormData(e.target);
  const entries = Object.fromEntries(formData.entries())

  if('user_id' in entries && !isNaN(entries['user_id'])) {
    const params = new window.URLSearchParams(window.location.search)
    params.set("user_id", entries['user_id'])
    params.set("limit", entries['limit'] ?? 1)
    params.set("offset", entries['offset'] ?? 0)

    window.location.search = params.toString()
    await resetForm(target)
  } 
}

async function resetForm(target) {
  target.classList.remove('loading')
}

function main() {
  const form = document.querySelector(selectors.form)
  if (!form) {
    console.error("Could not find form. Exiting...")
    return
  }

  form.addEventListener('submit', handleFormSubmit)
}

window.onload = main 
