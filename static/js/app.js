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
                if (this.commentSearch && repo.comment) {
                    if (!repo.comment.toLowerCase().includes(this.commentSearch.toLowerCase())) {
                        return false;
                    }
                }

                // Фильтрация по ключевому слову, если в keywordSearch есть текст
                if (this.keywordSearch && repo.key_word) {
                    if (!repo.key_word.toLowerCase().includes(this.keywordSearch.toLowerCase())) {
                        return false;
                    }
                }

                return true; // Если поля поиска пустые, не фильтруем
            });
        }
    },
    methods: {
        async fetchRepos() {
            this.isLoading = true;
            try {
                const response = await axios.get('/api/repos');
                this.repos = response.data;
            } catch (error) {
                console.error('Ошибка при получении данных:', error);
            } finally {
                this.isLoading = false;
            }
        },
        addComment(id, comment) {
            if (comment && comment.trim() !== '') {
                console.log(`Добавлен комментарий к репозиторию ${id}: ${comment}`);
                this.updateComment(id, comment);
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
                this.keyWord = ''; // очищаем поле ввода
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
                    alert("Должно помочь! Обнови страницу.");
                } else {
                    alert("Не удалось исправить базу данных. Попробуйте позже.");
                }
            } catch (error) {
                alert("Не помогло ((( Попробуй еще раз.");
                console.error('Ошибка при отправке запроса:', error);
            }
        }
    }
});
