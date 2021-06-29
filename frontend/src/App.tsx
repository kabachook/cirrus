import React from 'react';
import { Grommet, Main, Spinner, Text, DataTable, Grid, Box, Heading, Button, InfiniteScroll } from 'grommet';
import { Camera } from 'grommet-icons';
import useSWR from 'swr';
import { fetcher } from './fetcher';
import { API_URL } from './config';
import { useState } from 'react';

type Endpoint = {
  cloud: string,
  type: string,
  name: string,
  ip: string
}

type Snapshot = {
  timestamp: number,
  endpoints: Endpoint[]
}

function useSnapshots() {
  const { data, error } = useSWR<Snapshot[]>(API_URL + '/v1/snapshots', fetcher);

  return {
    snapshots: data?.sort((a, b) => b.timestamp - a.timestamp),
    isLoading: !error && !data,
    error
  }
}


function Table({ endpoints }: {
  endpoints?: Endpoint[]
}) {
  const { data, error } = useSWR(API_URL + '/v1/all', fetcher);


  if (!data) {
    return <Spinner />
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
      <DataTable columns={columns} data={endpoints?.length ? endpoints : data} />
    </React.Fragment>
  );
}

function Cirrus() {
  const { snapshots, isLoading, error } = useSnapshots();
  const [selectedEndpoints, setSelectedEndpoints] = useState<Endpoint[]>([])
  const [selectedTimestamp, setSelectedTimestamp] = useState(0)
  const [notification, setNotification] = useState('');

  if (isLoading) return <Spinner size="large" />
  if (error) return <Text>{error}</Text>


  const selectSnapshot = (ts: number) => {
    setSelectedTimestamp(ts);
    setSelectedEndpoints(snapshots!.find((snap) => snap.timestamp === ts)!.endpoints)
  }

  const newSnapshot = () => {
    fetcher(API_URL + "/v1/snapshot/new", { method: "POST" })
      .then(resp => {
        setNotification('Snapshot created');
        setTimeout(() => setNotification(''), 3000)
      })
  }

  return (
    <Grid
      rows={["xsmall", "auto"]}
      columns={["1/4", "3/4"]}
      gap="medium"
      areas={[
        { name: "sidebar_title", start: [0, 0], end: [0, 0] },
        { name: "table_title", start: [1, 0], end: [1, 0] },
        { name: "sidebar", start: [0, 1], end: [0, 1] },
        { name: "table", start: [1, 1], end: [1, 1] },
      ]}
    >

      <Box gridArea="sidebar_title">
        <Heading>
          Snapshots
        </Heading>
      </Box>
      <Box gridArea="table_title">
        <Heading>
          Cloud resources
        </Heading>
      </Box>
      <Box gridArea="sidebar" >
        {notification ?
          <Box>
            <Text>{notification}</Text>
          </Box> : null}
        <Button onClick={() => newSnapshot}>
          <Camera size="large" />
        </Button>
        <Box gap="medium">
          <InfiniteScroll items={snapshots}>
            {(snap: Snapshot) =>
              <Box
                flex={false}
                id={snap.timestamp.toString()}
                background={selectedTimestamp === snap.timestamp ? `neutral-3` : ''}
                gap="large"
                round="small"
              >
                <Text
                  onClick={() => selectSnapshot(snap.timestamp)}
                >
                  {new Date(snap.timestamp * 1000).toLocaleString()}
                </Text>
              </Box>
            }
          </InfiniteScroll>
        </Box>
      </Box>
      <Box gridArea="table" >
        <Table endpoints={selectedEndpoints} />
      </Box>
    </Grid>
  )
}

function App() {


  return (
    <Grommet plain>
      <Main pad="large">
        <Cirrus />
      </Main>
    </Grommet>
  )
}

export default App;
