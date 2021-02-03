import { extend } from 'umi-request';

const request = extend({
    prefix: 'http://localhost:7000',
    // timeout: 1000,
    // headers: {
    //     'Content-Type': 'multipart/form-data',
    // },
});


export async function getReviews(query: any) {
    const r = await request.get('/reviews', {
        params: query,
    })
        .then(function (resp) {
            return resp;
        })
        .catch(function (err) {
            console.error(err);
        });

    return r;
}

export async function getItems(query: any) {
    const r = await request.get('/items', {
        params: query,
    })
        .then(function (resp) {
            return resp;
        })
        .catch(function (err) {
            console.error(err);
        });

    return r;
}

export function startScrape(args: any) {
    request.post('/scrape', {
        data: args,
        requestType: "form",
    })
        .then(function (resp) {
            return resp;
        })
        .catch(function (err) {
            console.error(err);
        });
}