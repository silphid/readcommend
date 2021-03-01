import React  from "react";
import Criteria from "./Criteria"
import Results from "./Results"
import Container from '@material-ui/core/Container';
import Box from '@material-ui/core/Box';
import AppBar from '@material-ui/core/AppBar';

const App: React.FC = () => {
  return (
    <Container maxWidth="sm">
      <Box my={4}>
        <AppBar color="primary" position="static">
          <h1>Readcommend</h1>
        </AppBar>
        <Criteria />
        <Results />
      </Box>
    </Container>
  );
};

export default App;