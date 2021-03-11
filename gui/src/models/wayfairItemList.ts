import { wayfairCancel, wayfairGetItems, wayfairScrape, wayfairGetScrapeStatus } from '@/services/api';
import { message } from 'antd'
import { Reducer, Effect } from 'umi';

const delay = (ms: number) => new Promise(resolve => setTimeout(resolve, ms));

export interface IWayfairItemListModelState {
    items: any[],
    totalCount: number,
    pageSize: number,
    status: string,
    sku: string,
    itemNo: string,
    scrapeStatus: any,
    running: boolean,
}

export interface IItemListModel {
    // namespace: 'wayfairItemList';
    state: IWayfairItemListModelState;
    effects: {
        getItems: Effect;
        scrape: Effect;
        cancel: Effect;
        getScrapeStatus: Effect;
    };
    reducers: {
        setState: Reducer<IWayfairItemListModelState>;
    };
}

const ItemListModel: IItemListModel = {
    // namespace: 'wayfairItemList',
    state: {
        items: [],
        totalCount: 0,
        pageSize: 20,
        status: "",
        sku: "",
        itemNo: "",
        scrapeStatus: { Current: 0, TotalCount: 0 },
        running: false,
    },

    effects: {
        *scrape({ _ }, { call, put, select }) {
            const state = yield select((x: any) => x["wayfairItemList"]);
            const query = {
                status: state.status,
                sku: state.sku,
                itemNo: state.itemNo,
            };

            yield put({ type: 'setState', payload: { running: true } });
            const resp = yield call(wayfairScrape, query);
            console.log(resp);

            message.success("reviews scraping started");

            yield put({ type: 'getScrapeStatus' });
        },
        *getScrapeStatus({ _ }, { call, put, select }) {
            const resp = yield call(wayfairGetScrapeStatus);
            const running = resp.Current < resp.TotalCount;
            yield put({ type: 'setState', payload: { scrapeStatus: resp, running: running } });

            if (running) {
                yield delay(1000);
                yield put({ type: 'getScrapeStatus' });
            }
        },
        *getItems({ _ }, { call, put, select }) {
            const state = (yield select((x: any) => x["wayfairItemList"])) as IWayfairItemListModelState;
            const query = {
                status: state.status,
                sku: state.sku,
                itemNo: state.itemNo,
            };
            const resp = yield call(wayfairGetItems, query);
            yield put({ type: 'setState', payload: { items: resp.Items ?? [], totalCount: resp.TotalCount } });
        },
        *cancel({ _ }, { call, put, select }) {
            yield call(wayfairCancel);
            yield put({ type: 'setState', payload: { scrapeStatus: { Current: 0, TotalCount: 0 } } });
        },
    },
    reducers: {
        setState(state, action) {
            const r = {
                ...state,
                ...action.payload,
            };
            return r;
        },
    },
};

export default ItemListModel;