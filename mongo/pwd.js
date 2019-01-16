db = db.getSiblingDB('wisard')
db.createUser(
   {
     user:"mongoU",
     pwd:"password",
     roles:[ "dbOwner" ]
   }
);
