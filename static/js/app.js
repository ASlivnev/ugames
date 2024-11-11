new Vue({
    el: '#app',
    data: {
        repos: [],
        result: null,
        isLoading: false,
        keyWord: "",
        commentSearch: '',
        keywordSearch: ''
    },
    mounted() {
        this.fetchRepos();
    },
    computed: {
        filteredRepos() {
            return this.repos.filter(repo => {
                // Фильтрация по комментариям, если в commentSearch есть текст
                if (this.commentSearch) {
                    return repo.comment && repo.comment.toLowerCase().includes(this.commentSearch.toLowerCase());
                }

                if (this.keywordSearch) {
                    return repo.key_word && repo.key_word.toLowerCase().includes(this.keywordSearch.toLowerCase());
                }

                return true; // Если поле поиска пустое, не фильтруем
            });
        }
    },
    methods: {
        fetchRepos() {
            fetch('/api/repos')
                .then(response => response.json())
                .then(data => {
                    this.repos = data;
                })
                .catch(error => {
                    console.error('Ошибка при получении данных:', error);
                });
        },
        addComment(id, comment) {
            if (comment !== null) {
                console.log(`Добавлен комментарий к репозиторию ${id}: ${comment}`);
                this.updateComment(id, comment)
            }
        },
        async updateComment(id, comment) {
            try {
                const response = await axios.put('/api/addComment', {
                    id: id,
                    comment: comment
                });
                console.log('Comment updated successfully:', response.data);
                if (response.data.status !== 'Success') {
                    alert('Не удалось обновить комментарий! Попробуйте еще раз.');
                }

            } catch (error) {
                alert('Не удалось обновить комментарий! Попробуйте еще раз.');
                console.error('Error updating comment:', error);
            }
        },
        async submitForm() {
            this.isLoading = true;
            try {
                const response = await axios.post('/api/collectGitRepos', {
                    key_word: this.keyWord
                });
                this.result = response.data.message; // сохраняем результат в data
            } catch (error) {
                console.error('Ошибка при отправке запроса:', error);
                this.result = "Не удалось выполнить запрос."
            } finally {
                this.isLoading = false;
            }
        },
        async dbFix() {
            try {
                let response = await axios.get('/api/dbfix');
                if (response.data.status === 'Success') {
                    alert("Должно помочь! Обнови страницу.")
                }
            } catch (error) {
                alert("Не помогло ((( Попробуй еще раз.")
                console.error('Ошибка при отправке запроса:', error);
            }
        }
    }
});