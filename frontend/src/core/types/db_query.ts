import type {IQuery} from "./query.ts";

export interface IDBQuery {
    id: number;
    created_at: string;
    name: string;
    query: IQuery
}
