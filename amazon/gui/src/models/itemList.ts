import { getItems, startScrape } from '@/services/api';
import { message } from 'antd'
import { Reducer, Effect } from 'umi';

export interface IItemListModelState {
    items: any[],
    totalCount: number,
    pageSize: number,
    status: number,
    asin: string,
    itemNo: string,
}

export interface IItemListModel {
    namespace: 'itemList';
    state: IItemListModelState;
    effects: {
        getItems: Effect;
        loadMore: Effect;
        // search: Effect;
        scrape: Effect;
    };
    reducers: {
        setState: Reducer<IItemListModelState>;
    };
}

const ItemListModel: IItemListModel = {
    namespace: 'itemList',

    state: {
        items: [],
        totalCount: 0,
        pageSize: 10,
        status: -1,
        asin: "",
        itemNo: "",
    },

    effects: {
        // *search({ _ }, { call, put, select }) {
        //     const state = yield select((x: any) => x["itemList"]);
        //     const query = {
        //         pageSize: state.pageSize,
        //         status: state.status,
        //         asin: state.asin,
        //         itemNo: state.itemNo,
        //         searchAfter: "",
        //     };

        //     const resp = yield call(getItems, query);
        //     yield put({ type: 'setState', payload: { items: resp.Items, totalCount: resp.TotalCount } });
        // },
        *scrape({ _ }, { call, select }) {
            const state = yield select((x: any) => x["itemList"]);
            const query = {
                pageSize: state.pageSize,
                status: state.status,
                asin: state.asin,
                itemNo: state.itemNo,
                searchAfter: "",
            };

            const resp = yield call(startScrape, query);

            message.success(resp.count + " item(s) reviews scraped.");
        },
        *getItems({ _ }, { call, put, select }) {
            const state = (yield select((x: any) => x["itemList"])) as IItemListModelState;
            const query = {
                pageSize: state.pageSize,
                status: state.status,
                asin: state.asin,
                itemNo: state.itemNo,
                searchAfter: "",
            };
            const resp = yield call(getItems, query);
            yield put({ type: 'setState', payload: { items: resp.Items ?? [], totalCount: resp.TotalCount } });
        },
        *loadMore({ _ }, { call, put, select }) {
            const state = (yield select((x: any) => x["itemList"])) as IItemListModelState;
            const query = {
                pageSize: state.pageSize,
                status: state.status,
                asin: state.asin,
                itemNo: state.itemNo,
                searchAfter: "",
            };
            if (state.items.length > 0) {
                query.searchAfter = state.items[state.items.length - 1].SearchAfter;
            }
            const resp = yield call(getItems, query);
            const items = state.items.concat(resp.Items ?? []);
            yield put({ type: 'setState', payload: { items: items, totalCount: resp.TotalCount } });
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