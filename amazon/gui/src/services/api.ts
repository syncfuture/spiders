import { extend } from 'umi-request';
import u from '@/u'

const request = extend({
    prefix: u.BaseURI(),
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