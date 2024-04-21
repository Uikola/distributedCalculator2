// Обработчик отправки формы создания выражения
document.addEventListener('submit', function (event) {
    event.preventDefault();

    const expression = document.getElementById('expression').value;

    const expressionData = {
        expression: expression
    };

    authorizedFetch('http://localhost:8080/api/calculate', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(expressionData)
    })
        .then(response => {
            if (response.ok) {
                return response.json();
            } else if (response.status === 400) {
                errorMessage.textContent = 'Некорректное выражение';
            } else if (response.status === 500) {
                errorMessage.textContent = 'Ошибка на сервере или не хватает вычислительных ресурсов';
            } else {
                errorMessage.textContent = 'Неизвестная ошибка';
            }
        })
        .then(data => {
            // Добавляем выражение в список, если запрос успешен
            if (data) {
                addExpressionToList(data);

                updateExpressionList();
                // Сбрасываем значение поля ввода выражения
                document.getElementById('expression').value = '';
            }
        })
        .catch(error => {
            errorMessage.textContent = 'Ошибка на сервере';
        });
});

function addExpressionToList(expressionData) {
    const li = document.createElement('li');

    // Создаем элемент для отображения выражения
    const expressionSpan = document.createElement('span');
    expressionSpan.textContent = `${expressionData.expression} | `;
    li.appendChild(expressionSpan);

    // Создаем элемент для отображения статуса
    const statusSpan = document.createElement('span');
    statusSpan.textContent = `Статус: ${expressionData.status}`;
    statusSpan.classList.add(`status-${expressionData.status}`);
    li.appendChild(statusSpan);

    // Создаем элемент для отображения результата (если есть)
    if (expressionData.result !== undefined && expressionData.result !== null) {
        const resultSpan = document.createElement('span');
        resultSpan.textContent = ` | Результат: ${expressionData.result}`;
        li.appendChild(resultSpan);
    }

    expressionList.appendChild(li);
}

    // Получение списка выражений при загрузке страницы
    authorizedFetch('http://localhost:8080/api/expressions')
        .then(response => {
            if (response.ok) {
                if (response.status === 204) {
                    errorMessage.textContent = "У вас нет ни одного выражения";
                } else {
                    return response.json();
                }
            } else {
                errorMessage.textContent = 'Ошибка при получении списка выражений';
            }
        })
        .then(data => {
            // Добавляем каждое выражение из списка в HTML
            if (data) {
                data.forEach(addExpressionToList);
            }
        })
        .catch(error => {
            errorMessage.textContent = 'Ошибка на сервере';
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

// Функция для обновления списка выражений
function updateExpressionList() {
    authorizedFetch('http://localhost:8080/api/expressions')
        .then(response => {
            if (response.ok) {
                return response.json();
            } else {
                throw new Error('Ошибка при получении списка выражений');
            }
        })
        .then(data => {
            // Очищаем список выражений
            expressionList.innerHTML = '';
            // Добавляем каждое выражение из списка в HTML
            if (data) {
                data.forEach(addExpressionToList);
            }
        })
        .catch(error => {
            console.error(error);
        });
}

const logoutButton = document.getElementById('logoutButton');

// Добавляем обработчик события на кнопку "Выход"
logoutButton.addEventListener('click', function() {
    // Удаляем токен из локального хранилища
    localStorage.removeItem('token');
    // Перенаправляем пользователя на страницу входа
    window.location.href = '../html/login.html'; // Замените ссылку на свою страницу входа
});