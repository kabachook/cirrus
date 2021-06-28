import React from 'react';
import { Grommet, Main, Heading, Text, DataTable } from 'grommet';
import useSWR from 'swr';
import { fetcher } from './fetcher';
import { API_URL } from './config';

function Table() {
  const { data, error } = useSWR(API_URL + '/v1/all', fetcher);


  if (!data) {
    return <Heading>Loading...</Heading>
  }

  if (error) {
    return <Text>Error while loading: {error.toString()}</Text>
  }

  const columns = [
    {
      property: "cloud",
      header: "☁️"
    }, {
      property: "type",
      header: "Type"
    }, {
      property: "name",
      header: "Name"
    }, {
      property: "ip",
      header: "IP"
    }]

  return (
    <React.Fragment>
      <DataTable columns={columns} data={data} />
    </React.Fragment>
  );
}

function App() {
  return (
    <Grommet plain>
      <Main pad="large">
        <Table />
      </Main>
    </Grommet>
  )
}

export default App;
