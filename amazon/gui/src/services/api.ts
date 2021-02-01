import request from 'umi-request';

export async function getReviews() {
    const r = await request.get('http://localhost:7000/reviews')
        .then(function (resp) {
            return resp;
        })
        .catch(function (err) {
            console.error(err);
        });

    return r;
}

export async function getItems(args: any) {
    const r = await request.get('http://localhost:7000/items', {
        params: args,
    })
        .then(function (resp) {
            return resp;
        })
        .catch(function (err) {
            console.error(err);
        });

    return r;
}

export async function startScrape(args: any) {
    const r = await request.post('http://localhost:7000/scrape', {
        data: args,
    })
        .then(function (resp) {
            return resp;
        })
        .catch(function (err) {
            console.error(err);
        });

    return r;
}