const delimiters = ['${', '}'];

const main = {
    delimiters,
    data() {
        return {
            workers: [],
        };
    },
    mounted() {
        document.body.removeAttribute('hidden');
        document.addEventListener("click", async () => {
            await this.summon();
        });
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
            worker.delay = worker.delay || 0;
            this.setWorkers([...this.workers, worker]);
            await (() => new Promise((resolve) => setTimeout(() => resolve(), worker.delay * 1000)))();
            await this.putWorker(worker);
            await (() => new Promise((resolve) => setTimeout(() => resolve(), 500)))();
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

const progress = {
    delimiters,
    template: '#app-progress',
    props: {
        delay: {
            type: Number,
            required: true,
        },
    },
    data() {
        return {
            progress: 0,
        };
    },
    mounted() {
        const timer = setInterval(() => {
            this.progress++;
        }, 1000);
        setTimeout(() => {
            clearInterval(timer);
        }, this.delay * 1000);
    },
};

Vue
    .createApp(main)
    .component('app-progress', progress)
    .mount('#app');
