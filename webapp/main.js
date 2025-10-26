window.addEventListener('DOMContentLoaded', () => {
    const server_url = localStorage.getItem('server_url');
    const api_key = localStorage.getItem('api_key');
    if (!server_url || !api_key) {
        window.location.href = 'setup.html';
    }
    navigation();
    upload(server_url, api_key);
    history(server_url, api_key);
});

function navigation() {
    const nav_upload = document.getElementById('nav_upload');
    const nav_history = document.getElementById('nav_history');
    const upload = document.getElementById('upload');
    const history = document.getElementById('history');

    nav_upload.addEventListener('click', () => {
        upload.classList.remove('hide');
        history.classList.add('hide');
        nav_history.classList.remove('active');
        nav_upload.classList.add('active');
    });
    nav_history.addEventListener('click', () => {
        upload.classList.add('hide');
        history.classList.remove('hide');
        nav_upload.classList.remove('active');
        nav_history.classList.add('active');
    });
}

function upload(base_url = '/', api_key = '') {
    const name = document.getElementById('name');
    const file = document.getElementById('file');
    const progress = document.getElementById('progress');
    const progress_bar = document.getElementById('progress-bar');
    const upload_form = document.getElementById('upload_form');
    const status_bar = document.getElementById('status_bar');
    const last_link = document.getElementById('last_link');


    upload_form.addEventListener('submit', e => {
        e.preventDefault();

        const formData = new FormData();
        formData.append('name', name.value);
        formData.append('file', file.files[0]);

        const xhr = new XMLHttpRequest();
        xhr.upload.addEventListener('progress', e => {
            if (e.lengthComputable) {
                const percent = (e.loaded / e.total) * 100;
                if (status_bar.classList.contains('hide')) {
                    status_bar.classList.remove('hide');
                }
                progress_bar.style.width = `${percent}%`;
                progress.style.width = `${percent}%`;
            }
        });

        xhr.onload = () => {
            const response = JSON.parse(xhr.responseText);
            status_bar.classList.add('hide');
            if (response.status === 'success') {
                progress_bar.style.width = '100%';
                progress.style.width = '100%';
                name.value = '';
                file.value = '';
                last_link.innerHTML = `<a href="${response.url}">${response.url}</a>`;
            } else {
                progress_bar.style.width = '0%';
                progress.style.width = '0%';
            }
        };

        xhr.open('POST', base_url);
        xhr.setRequestHeader('Authorization', `Bearer ${api_key}`);
        xhr.send(formData);
    });
}

function history(base_url = '/', api_key = '') {
    const table = document.getElementById('history_table');
    const tbody = table.getElementsByTagName('tbody')[0];
    const last_link = document.getElementById('last_link');

    fetch(`${base_url}uploads`, {
        method: 'GET',
        headers: {
            'Authorization': `Bearer ${api_key}`,
        }
    })
        .then(response => response.json())
        .then(data => {
            if (data.status === 'success') {
                data.data.forEach(item => {
                    const row = document.createElement('tr');
                    row.innerHTML = `
                        <td><a href="${base_url}${item.short_code}">${item.short_code}</a></td>
                        <td>${item.filename}</td>
                    `;
                    tbody.appendChild(row);
                });
                last_data = data.data[0];
                last_link.innerHTML = `<a href="${base_url}${last_data.short_code}">${last_data.short_code}</a>`;
            }
        });
}