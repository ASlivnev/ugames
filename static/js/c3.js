new Vue({
    el: '#app',
    data: {
        games: [],
        result: null,
        isLoading: false,
        keyWord: "",
        commentSearch: '',
        keywordSearch: ''
    },
    mounted() {
        this.fetchGames("all");
    },
    computed: {

    },
    methods: {
        async fetchGames(filter) {
            this.isLoading = true;
            try {
                const response = await axios.get('/api/construct/' + filter);
                this.games = response.data;
            } catch (error) {
                console.error('Ошибка при получении данных:', error);
            } finally {
                this.isLoading = false;
            }
        },
        async saveGame(game) {
            try {
                const response = await axios.put("/api/updateC3Game", {
                    id: game.id,
                    list: game.list,
                    comment: game.comment,
                });

                // Обработка успешного ответа
                console.log(`Игра ${game.id} сохранена:`, response.data);
                alert(`Игра ${game.id} успешно сохранена!`);
            } catch (error) {
                // Обработка ошибки
                console.error(`Ошибка при сохранении игры ${game.id}:`, error);
                alert(`Ошибка при сохранении игры ${game.id}.`);
            }
        }



    }
});
