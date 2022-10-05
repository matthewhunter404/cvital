
export function getJwtToken() {
    return localStorage.getItem('jwt')
}

export function setJwtToken(token) {
    localStorage.setItem('jwt', token)
}
