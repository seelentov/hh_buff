import './App.css'
import {useEffect, useState} from "react";
import {CartesianGrid, Legend, Line, LineChart, ResponsiveContainer, Tooltip, XAxis, YAxis} from "recharts";

interface IQuery {
    text: string;
}

interface IDBQuery {
    id: number;
    created_at: string;
    name: string;
    query: IQuery
}

interface IDBSnaphot {
    id: number;
    created_at: string;
    count: number;
    query_id: number;
    query: IDBQuery
}

function App() {
    const [data, setData] = useState<IDBSnaphot[]>([]);
    const [isLoading, setIsLoading] = useState<boolean>(true);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const response = await fetch('/rest/data');
                if (!response.ok) {
                    throw new Error(`HTTP error! status: ${response.status}`);
                }
                const result = await response.json();
                // console.log(result);
                setData(result);
            } catch (err) {
                console.error(err);
            } finally {
                setIsLoading(false);
            }
        };

        fetchData();
    }, []);

    const prepareChartData = () => {
        const groups: Record<string, any> = {};

        data.forEach(item => {
            const time = new Date(item.created_at).toLocaleTimeString(); // Время для оси X
            if (!groups[time]) {
                groups[time] = {time};
            }
            const lineName = item.query.name || `Query ${item.query_id}`;
            groups[time][lineName] = item.count;
        });

        return Object.values(groups).sort((a, b) =>
            new Date(a.time).getTime() - new Date(b.time).getTime()
        );
    };

    const chartData = prepareChartData();

    const queryNames = Array.from(new Set(data.map(item => item.query.name || `Query ${item.query_id}`)));
    const colors = ['#8884d8', '#82ca9d', '#ffc658', '#ff7300'];

    return (
        <main style={{width: '100%', height: 400}}>
            {isLoading ? (
                <p>Loading...</p>
            ) : (
                <ResponsiveContainer width="100%" height="100%">
                    <LineChart data={chartData}>
                        <CartesianGrid
                            strokeDasharray="3 3"
                        />
                        <XAxis
                            dataKey="time"
                        />
                        <YAxis
                            domain={[0, 'dataMax']}
                            tickFormatter={(val) => val.toLocaleString()}
                        />
                        <Tooltip/>
                        <Legend/>
                        {queryNames.map((name, index) => (
                            <Line
                                key={name}
                                type="monotone"
                                dataKey={name}
                                stroke={colors[index % colors.length]}
                                dot={false}
                            />
                        ))}
                    </LineChart>
                </ResponsiveContainer>
            )}
        </main>
    );
}

export default App
