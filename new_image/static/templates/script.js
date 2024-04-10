range=document.getElementById('myRange');
val=document.getElementById('value');
val.value="Value: "+range.value;


  range.addEventListener('change', function(){
    val.value="Value: "+range.value;
    
  })