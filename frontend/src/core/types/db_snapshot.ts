import type {IDBQuery} from "./db_query.ts";

export interface IDBSnaphot {
    id: number;
    created_at: string;
    count: number;
    query_id: number;
    query: IDBQuery
}