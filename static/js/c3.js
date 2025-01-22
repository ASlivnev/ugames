new Vue({
    el: '#app',
    data: {
        games: [],
        result: null,
        isLoading: false,
        keyWord: "",
        commentSearch: '',
        keywordSearch: '',
        search: '',
    },
    mounted() {
        this.fetchGames("all");
    },
    computed: {

    },
    methods: {
        async fetchSearchGames() {
            this.isLoading = true;
            try {
                const response = await axios.post("/api/searchC3Games", {
                    searchRequest: this.search
                });
                this.games = response.data;
            } catch (error) {
                console.error('Ошибка при получении данных:', error);
            } finally {
                this.isLoading = false;
            }
        },
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
                this.$set(game, 'hidden', true);
                //alert(`Игра ${game.id} успешно сохранена!`);
            } catch (error) {
                // Обработка ошибки
                console.error(`Ошибка при сохранении игры ${game.id}:`, error);
                alert(`Ошибка при сохранении игры ${game.id}.`);
            }
        },
        async saveGreyGame(game) {
            try {
                const response = await axios.put("/api/updateC3Game", {
                    id: game.id,
                    list: "grey",
                    comment: game.comment,
                });

                // Обработка успешного ответа
                console.log(`Игра ${game.id} сохранена:`, response.data);
                this.$set(game, 'hidden', true);
                //alert(`Игра ${game.id} успешно сохранена!`);
            } catch (error) {
                // Обработка ошибки
                console.error(`Ошибка при сохранении игры ${game.id}:`, error);
                alert(`Ошибка при сохранении игры ${game.id}.`);
            }
        }




    }
});
