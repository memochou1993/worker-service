const app = {
    delimiters: ['${', '}'],
    data() {
        return {
            workers: [],
        };
    },
    mounted() {
        document.body.removeAttribute('hidden');
    },
    methods: {
        setWorkers(workers) {
            this.workers = workers;
        },
        async summon() {
            const { worker } = await this.fetchWorker();
            if (!worker) {
                return;
            }
            this.setWorkers([...this.workers, worker]);
            await (() => new Promise((resolve) => setTimeout(() => resolve(), worker.delay * 1000)))();
            await this.putWorker(worker);
            this.setWorkers(this.workers.filter(w => w.number !== worker.number));
        },
        fetchWorker() {
            return fetch('api/worker')
                .then((res) => res.ok ? res.json() : Promise.reject(res.statusText))
                .then((res) => res)
                .catch(() => Object);
        },
        putWorker(worker) {
            const init = {
                body: JSON.stringify({ number: worker.number }),
                headers: { 'content-type': 'application/json' },
                method: 'PUT',
            };
            return fetch('api/worker', init)
                .then();
        },
    },
};

Vue.createApp(app).mount('#app');
