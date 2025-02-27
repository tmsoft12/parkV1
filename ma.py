from sklearn.linear_model import LinearRegression
import numpy as np


X= np.array([[50,1],[70,2],[100,3],[120,4],[150,5]])
y=np.array([15.0000,20.0000,25.0000,90.0000,35.0000])


model =LinearRegression()

model.fit(X,y)


height = float(input("cm"))
weight  = float(input("kg"))

new_body= np.array([[height,weight]])
prediction = model.predict(new_body)
print(f"{prediction[0]}")