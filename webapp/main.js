window.addEventListener('DOMContentLoaded', async () => {
    const server_url = await getValue('server_url');
    if (!server_url) {
        window.location.href = '/setup';
        return;
    }

    const isSessionValid = await checkSession(server_url);
    if (!isSessionValid) {
        window.location.href = '/setup';
        return;
    }

    navigation();
    upload(server_url);
    await history(server_url);
});

async function checkSession(base_url) {
    try {
        const response = await fetch(`${base_url}session/validate`, {
            method: 'GET',
            credentials: 'include'
        });

        if (!response.ok) {
            return false;
        }

        const data = await response.json();
        return data.status === 'success';
    } catch (err) {
        console.error('Session validation failed:', err);
        return false;
    }
}

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

function upload(base_url = '/') {
    const file = document.getElementById('file');
    const progress = document.getElementById('progress');
    const progress_bar = document.getElementById('progress-bar');
    const upload_form = document.getElementById('upload_form');
    const status_bar = document.getElementById('status_bar');
    const last_link = document.getElementById('last_link');
    const file_label_text = document.getElementById('file_label_text');

    file.onchange = () => {
        file_label_text.innerHTML = file.files[0].name;
    };

    upload_form.addEventListener('submit', e => {
        e.preventDefault();

        const formData = new FormData();
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
            if (xhr.status === 401) {
                window.location.href = '/setup';
                return;
            }

            const response = JSON.parse(xhr.responseText);
            status_bar.classList.add('hide');

            if (response.status === 'success') {
                progress_bar.style.width = '100%';
                progress.style.width = '100%';
                file.value = '';
                last_link.innerHTML = `<a href="${response.url}">${response.url}</a>`;
                file_label_text.innerHTML = 'select file to upload';
            } else {
                progress_bar.style.width = '0%';
                progress.style.width = '0%';
            }
        };

        xhr.open('POST', `${base_url}upload`);
        xhr.withCredentials = true;
        xhr.send(formData);
    });
}

async function history(base_url = '/') {
    const table = document.getElementById('history_table');
    const tbody = table.getElementsByTagName('tbody')[0];
    const last_link = document.getElementById('last_link');

    try {
        const response = await fetch(`${base_url}uploads`, {
            method: 'GET',
            credentials: 'include'
        });

        if (response.status === 401) {
            window.location.href = '/setup';
            return;
        }

        const data = await response.json();
        if (!data || data.status !== 'success') return;

        tbody.innerHTML = '';

        for (const item of data.data) {
            const row = document.createElement('tr');
            row.innerHTML = `
                <td title="${item.filename}">${item.filename}</td>
                <td><a href="${base_url}${item.short_code}">${item.short_code}</a></td>
            `;
            tbody.appendChild(row);
        }

        const last_data = data.data[0];
        if (last_data) {
            last_link.innerHTML = `<a href="${base_url}${last_data.short_code}">${base_url}${last_data.short_code}</a>`;
        }
    } catch (err) {
        console.error('Failed to load history:', err);
    }
}
