import React  from "react";
import Criteria from "./criteria/Criteria"
import Results from "./Results"
import Container from '@material-ui/core/Container';
import Box from '@material-ui/core/Box';
import AppBar from '@material-ui/core/AppBar';
import * as models from './criteria/models';

const App: React.FC = () => {
  const [criteria, setCriteria] = React.useState<models.Criteria>(models.DefaultCriteria);

  return (
    <Container maxWidth="md">
      <Box my={4}>
        <AppBar color="primary" position="static">
          <h1>Readcommend</h1>
        </AppBar>
        <Criteria onChange={(x) => setCriteria(x)} />
        <Results />
        <p>Authors: {criteria.authors.join(' ')}</p>
        <p>Genres: {criteria.genres.join(' ')}</p>
        <p>Era minYear: {criteria.minYear} maxYear: {criteria.maxYear}</p>
        <p>Size minPages: {criteria.minPages} maxPages: {criteria.maxPages}</p>
        <p>Limit {criteria.limit}</p>
      </Box>
    </Container>
  );
};

export default App;