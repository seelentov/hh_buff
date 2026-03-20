import './App.css'
import {useEffect, useState} from "react";
import Chart from "./components/Chart.tsx";
import type {IDBSnaphot} from "./core/types/db_snapshot.ts";
import Select, {type MultiValue} from 'react-select';
import DatePicker from "react-datepicker";

interface IOption {
    value: string;
    label: string;
}

const options: IOption[] = [
    {value: 'chocolate', label: 'Chocolate'},
    {value: 'strawberry', label: 'Strawberry'},
    {value: 'vanilla', label: 'Vanilla'},
];

function App() {
    const [data, setData] = useState<IDBSnaphot[]>([]);
    const [isLoading, setIsLoading] = useState<boolean>(true);

    const [selectedQueries, setSelectedQueries] = useState<IOption[]>([]);

    const [startDate, setStartDate] = useState<Date | null>(new Date());
    const [endDate, setEndDate] = useState<Date | null>(new Date());

    useEffect(() => {
        const fetchData = async () => {
            try {
                const response = await fetch('/rest/data');
                if (!response.ok) {
                    throw new Error(`HTTP error! status: ${response.status}`);
                }
                const result = await response.json();
                setData(result);
            } catch (err) {
                console.error(err);
            } finally {
                setIsLoading(false);
            }
        };

        fetchData();
    }, []);

    const handleChange = (newValue: MultiValue<IOption>) => {
        setSelectedQueries(newValue as IOption[]);
    };

    return (
        <div className="wrapper">
            <main>
                <section id="chart">
                    <div className="container-v">
                        <Chart data={data} isLoading={isLoading}/>
                    </div>
                    <div className="container container-v flex gap">
                        <div style={{width: '100%'}}>
                            <Select
                                options={options}
                                isMulti
                                value={selectedQueries}
                                onChange={handleChange}
                                isSearchable
                            />
                        </div>
                        <DatePicker
                            selected={startDate}
                            onChange={(date: Date | null) => setStartDate(date)}
                            locale="ru"
                            dateFormat="dd.MM.yyyy"
                            className="form-control"
                            placeholderText="Нажмите, чтобы выбрать"
                            isClearable
                            showPopperArrow={false} // Убираем стрелочку сверху для минимализма
                        />
                        <DatePicker
                            selected={endDate}
                            onChange={(date: Date | null) => setEndDate(date)}
                            locale="ru"
                            dateFormat="dd.MM.yyyy"
                            className="form-control"
                            placeholderText="Нажмите, чтобы выбрать"
                            isClearable
                            showPopperArrow={false}
                        />
                    </div>
                    <div className="container">
                        <button type="button" className="btn btn-primary">Обновить график</button>
                    </div>
                </section>
            </main>
        </div>

    );
}

export default App
