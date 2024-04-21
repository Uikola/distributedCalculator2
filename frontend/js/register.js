document.addEventListener('DOMContentLoaded', function () {
    const registerForm = document.getElementById('registerForm');
    const errorMessage = document.getElementById('errorMessage');

    registerForm.addEventListener('submit', function (event) {
        event.preventDefault();

        const login = document.getElementById('login').value;
        const password = document.getElementById('password').value;

        const userData = {
            login: login,
            password: password
        };

        fetch('http://localhost:8080/api/register', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(userData)
        })
            .then(response => {
                if (response.ok) {
                    const form = document.createElement('form');
                    form.method = 'POST';
                    form.action = '../html/login.html';
                    document.body.appendChild(form);
                    form.submit();
                } else if (response.status === 400) {
                    errorMessage.textContent = 'Неверный запрос или JSON';
                } else if (response.status === 409) {
                    errorMessage.textContent = 'Пользователь с таким логином уже существует';
                } else {
                    errorMessage.textContent = 'Ошибка на сервере';
                }
            })
            .catch(error => {
                errorMessage.textContent = 'Ошибка на сервере';
            });
    });
});
