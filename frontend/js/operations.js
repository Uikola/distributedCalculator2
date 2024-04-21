document.addEventListener('DOMContentLoaded', function () {
    const operationForm = document.getElementById('operationForm');
    const operationList = document.getElementById('operationList');
    const errorMessage = document.getElementById('errorMessage');

    // Получение списка операций с сервера и заполнение выпадающего списка
    function fetchOperation() {
        authorizedFetch('http://localhost:8080/api/operations')
            .then(response => {
                if (response.ok) {
                    return response.json();
                } else {
                    throw new Error('Ошибка при получении списка операций');
                }
            })
            .then(data => {
                const operationSelect = document.getElementById('operation');
                // Очищаем выпадающий список перед добавлением новых операций
                operationSelect.innerHTML = '';

                // Добавляем каждую операцию в выпадающий список
                for (const operation in data) {
                    const option = document.createElement('option');
                    option.value = operation;
                    option.textContent = operation;
                    operationSelect.appendChild(option);
                }

                // Отображаем список операций на странице
                displayOperations(data);
            })
            .catch(error => {
                console.error('Произошла ошибка:', error);
            });
    }
    fetchOperation()

    // Обработчик отправки формы изменения операции
    operationForm.addEventListener('submit', function (event) {
        event.preventDefault();

        const selectedOperation = document.getElementById('operation').value;
        const time = parseInt(document.getElementById('time').value, 10);

        const operationData = {
            operation: selectedOperation,
            time: time
        };

        authorizedFetch('http://localhost:8080/api/operations', {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(operationData)
        })
            .then(response => {
                if (response.ok) {
                    fetchOperation()
                    return response.json();
                } else {
                    throw new Error('Ошибка при изменении операции');
                }
            })
            .then(data => {
                // Обновляем список операций после изменения
                operationList.textContent = '';
                displayOperations(data);
            })
            .catch(error => {
                errorMessage.textContent = error.message;
            });
    });

    // Функция для отображения списка операций
    function displayOperations(operations) {
        operationList.innerHTML = '';
        for (const operation in operations) {
            const listItem = document.createElement('li');
            listItem.textContent = `${operation}: ${operations[operation]} секунд`;
            operationList.appendChild(listItem);
        }
    }
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