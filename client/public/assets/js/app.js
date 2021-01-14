const delimiters = ['${', '}'];

const main = {
    delimiters,
    data() {
        return {
            workers: [],
            summoned: 0,
            gems: 0,
        };
    },
    mounted() {
        this.initialize();
        document.body.removeAttribute('hidden');
        window.onbeforeunload = async () => {};
    },
    methods: {
        setWorkers(workers) {
            this.workers = workers;
        },
        setSummoned(summoned) {
            this.summoned = summoned;
        },
        setGems(gems) {
            this.gems = gems;
        },
        async initialize() {
            const numbers = Array(30).fill(0).map((_, i) => i + 1);
            this.changeCursor('progress');
            await Promise.all(numbers.map(() => this.fetchWorker()));
            await Promise.all(numbers.sort(() => Math.random() - 0.5).map((n) => this.putWorker(n)));
            this.changeCursor('grab');
            document.addEventListener('click', () => this.summon());
        },
        async summon() {
            const { worker } = await this.fetchWorker();
            if (!worker) {
                return;
            }
            worker.delay = worker.delay || 0;
            this.setWorkers([...this.workers, worker]);
            this.setSummoned(this.summoned + 1);
            await this.delay(worker.delay * 1000 + 250);
            await this.putWorker(worker.number);
            this.setWorkers(this.workers.filter(w => w.number !== worker.number));
            this.setGems(this.gems + worker.delay);
        },
        fetchWorker() {
            return fetch('api/worker')
                .then((res) => res.ok ? res.json() : Promise.reject(res.statusText))
                .then((res) => res)
                .catch(async () => {
                    await this.changeCursor('progress');
                    await this.delay(1000);
                    await this.changeCursor('grab');
                    return Object;
                });
        },
        putWorker(number) {
            return fetch('api/worker', {
                body: JSON.stringify({ number }),
                headers: { 'content-type': 'application/json' },
                method: 'PUT',
            })
                .then()
                .catch((err) => {
                    console.log(err);
                });
        },
        changeCursor(cursor) {
            document.querySelector('html').style.cursor = cursor;
        },
        getImage(number, delay) {
            if (delay === 10) {
                return '9';
            }
            switch (true) {
                case this.gems >= 5000:
                    return `${number % 8 + 1}-${5}`;
                case this.gems >= 2500 && delay >= 8:
                    return `${number % 8 + 1}-${5}`;
                case this.gems >= 2500:
                    return `${number % 8 + 1}-${4}`;
                case this.gems >= 1000 && delay >= 8:
                    return `${number % 8 + 1}-${4}`;
                case this.gems >= 1000:
                    return `${number % 8 + 1}-${3}`;
                case this.gems >= 500 && delay >= 8:
                    return `${number % 8 + 1}-${3}`;
                case this.gems >= 500:
                    return `${number % 8 + 1}-${2}`;
                case this.gems >= 100 && delay >= 8:
                    return `${number % 8 + 1}-${2}`;
                default:
                    return `${number % 8 + 1}-${1}`;
            }
        },
        delay(milliseconds) {
            return new Promise((resolve) => setTimeout(() => resolve(), milliseconds));
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
