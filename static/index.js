const selectors = {
  form: "#form--name",
  quote: "#quote"
}

async function handleFormSubmit(e) {
  e.preventDefault()
  e.stopPropagation()

  const { target } = e;

  if (target && !target.classList.contains("loading")) {
    target.classList.add("loading")
  }

  const response = await fetch("/api/kanye")
  const dataOrError = await response.json()
  if (!response.ok || response.status < 200 || response.status > 299) {
    console.warn("API request failed.")
    console.error(dataOrError)
    await resetForm(target)
    alert("Couldn't get your quote. Sorry about that. Please try again later.")
    return
  } 

  await render(dataOrError)
  await resetForm(target)
}

async function render(data) {
  const quote = document.querySelector(selectors.quote)
  if(!quote){
    console.error("Could not find quote. Exiting...")
    return
  }

  quote.innerHTML = data.message
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
