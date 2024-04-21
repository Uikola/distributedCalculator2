document.addEventListener('DOMContentLoaded', function () {
    const loginForm = document.getElementById('loginForm');
    const errorMessage = document.getElementById('errorMessage');

    loginForm.addEventListener('submit', function (event) {
        event.preventDefault();

        const login = document.getElementById('login').value;
        const password = document.getElementById('password').value;

        const userData = {
            login: login,
            password: password
        };

        fetch('http://localhost:8080/api/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(userData)
        })
            .then(response => {
                if (response.ok) {
                    return response.json();
                } else if (response.status === 400) {
                    errorMessage.textContent = 'Неверный запрос или JSON';
                } else if (response.status === 404) {
                    errorMessage.textContent = 'Пользователь с таким логином и паролем не найден';
                } else {
                    errorMessage.textContent = 'Ошибка на сервере';
                }
            })
            .then(data => {
                if (data) {
                    const token = data.token;
                    // Сохраняем токен в localStorage
                    localStorage.setItem('token', token);
                    // Перенаправляем на главную страницу приложения
                    window.location.href = '../html/main.html';
                }
            })
            .catch(error => {
                errorMessage.textContent = 'Ошибка на сервере';
            });
    });
});
