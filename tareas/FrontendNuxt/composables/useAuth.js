export const useAuth = () => {
  const userId = ref(null)

  //Description: This composable handles user authentication, including login and logout functionality.
  const logout = () => {
    localStorage.removeItem('token')
    localStorage.removeItem('userId') 
    navigateTo('/login')
  }

  // load userId from localStorage on mount
  onMounted(() => {
    const storedUserId = localStorage.getItem('userId')
    if (storedUserId) {
      userId.value = storedUserId
    }
  })

  return {
    userId, 
    logout
  }
}
