document.addEventListener('DOMContentLoaded', function () {
    const resourceList = document.getElementById('resourceList');

    // Отправка запроса на получение данных о вычислительных ресурсах
    authorizedFetch('http://localhost:8080/api/cresources')
        .then(response => {
            if (response.ok) {
                return response.json();
            } else {
                throw new Error('Ошибка при получении данных о вычислительных ресурсах');
            }
        })
        .then(data => {
            // Добавление элементов списка вычислительных ресурсов
            for (const [resourceName, expression] of Object.entries(data)) {
                const listItem = document.createElement('li');
                listItem.textContent = `${resourceName}: ${expression || 'свободен'}`;
                resourceList.appendChild(listItem);
            }
        })
        .catch(error => {
            console.error(error);
            // Обработка ошибок, например, отображение сообщения пользователю
        });
});

function authorizedFetch(url, options) {
    // Проверяем, определен ли options
    if (!options) {
        options = {};
    }

    // Проверяем, определен ли options.headers
    if (!options.headers) {
        options.headers = {};
    }

    // Извлекаем токен из localStorage
    const token = localStorage.getItem('token');
    // Добавляем токен в заголовок Authorization
    options.headers = {
        ...options.headers,
        'Authorization': `Bearer ${token}`
    };

    // Отправляем запрос с обновленными заголовками
    return fetch(url, options);
}

const logoutButton = document.getElementById('logoutButton');

// Добавляем обработчик события на кнопку "Выход"
logoutButton.addEventListener('click', function() {
    // Удаляем токен из локального хранилища
    localStorage.removeItem('token');
    // Перенаправляем пользователя на страницу входа
    window.location.href = '../html/login.html'; // Замените ссылку на свою страницу входа
});
