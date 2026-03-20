import {CartesianGrid, Line, LineChart, ResponsiveContainer, Tooltip, XAxis, YAxis} from "recharts";
import type {IDBSnaphot} from "../core/types/db_snapshot.ts";
import {useMemo} from "react";
import {ThreeDots} from 'react-loader-spinner';
import './Chart.css'

interface IChartData {
    data: IDBSnaphot[];
    isLoading?: boolean;
}

function Chart({data, isLoading}: IChartData) {
    const {chartData, queries} = useMemo(() => {
        const groups: Record<string, any> = {};
        const uniqueQueries = new Map<number, string>();

        data.forEach(item => {
            const date = new Date(item.created_at).toLocaleDateString(); // Или другой формат даты
            const queryName = item.query.name;
            const queryId = item.query_id;

            if (!groups[date]) {
                groups[date] = {date};
            }
            groups[date][queryName] = item.count;
            uniqueQueries.set(queryId, queryName);
        });

        return {
            chartData: Object.values(groups),
            queries: Array.from(uniqueQueries.values())
        };
    }, [data]);

    if (isLoading) {
        return (
            <div style={{display: 'flex', justifyContent: 'center', alignItems: 'center', height: '600px'}}>
                <ThreeDots height="80" width="80" color="#0b69ff" ariaLabel="loading"/>
            </div>
        );
    }

    return <ResponsiveContainer width="100%" height={600}>
        <LineChart data={chartData}>
            <CartesianGrid strokeDasharray="3 3" vertical={false}/>
            <XAxis
                dataKey="date"
                minTickGap={60}
            />
            <YAxis/>
            <Tooltip/>
            {queries.map((name, index) => (
                <Line
                    key={name}
                    type="monotone"
                    dataKey={name}
                    stroke={getColor(index)}
                    dot={false}
                    strokeWidth={1}
                />
            ))}
        </LineChart>
    </ResponsiveContainer>
}

const getColor = (index: number) => {
    const colors = ['#f1c40f', '#e74c3c', '#3498db', '#2c3e50', '#9b59b6', '#e67e22'];
    return colors[index % colors.length];
};

export default Chart