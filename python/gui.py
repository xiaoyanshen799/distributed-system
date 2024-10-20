import tkinter as tk
from tkinter import messagebox, filedialog, ttk
import grpc
import pet_pb2
import pet_pb2_grpc


class PetApp:
    def __init__(self, root):
        self.root = root
        self.root.title("Pet Registration and Search")

        self.create_widgets()

        # Initialize gRPC client
        self.channel = grpc.insecure_channel('localhost:50051')
        self.stub = pet_pb2_grpc.PetServiceStub(self.channel)

    def create_widgets(self):
        # Registration Frame
        self.registration_frame = tk.Frame(self.root)
        self.registration_frame.pack(pady=10)

        tk.Label(self.registration_frame, text="Pet Name:").grid(row=0, column=0)
        self.pet_name_entry = tk.Entry(self.registration_frame)
        self.pet_name_entry.grid(row=0, column=1)

        tk.Label(self.registration_frame, text="Gender:").grid(row=1, column=0)
        self.pet_gender_entry = tk.Entry(self.registration_frame)
        self.pet_gender_entry.grid(row=1, column=1)

        tk.Label(self.registration_frame, text="Age:").grid(row=2, column=0)
        self.pet_age_entry = tk.Entry(self.registration_frame)
        self.pet_age_entry.grid(row=2, column=1)

        tk.Label(self.registration_frame, text="Breed:").grid(row=3, column=0)
        self.pet_breed_entry = tk.Entry(self.registration_frame)
        self.pet_breed_entry.grid(row=3, column=1)

        tk.Label(self.registration_frame, text="Photo:").grid(row=4, column=0)
        self.pet_photo_entry = tk.Entry(self.registration_frame)
        self.pet_photo_entry.grid(row=4, column=1)
        tk.Button(self.registration_frame, text="Browse", command=self.browse_photo).grid(row=4, column=2)

        tk.Button(self.registration_frame, text="Register Pet", command=self.register_pet).grid(row=5, column=0,
                                                                                                columnspan=3)

        # Search Frame
        self.search_frame = tk.Frame(self.root)
        self.search_frame.pack(pady=10)

        tk.Label(self.search_frame, text="Select Search Field:").grid(row=0, column=0)
        self.search_field = ttk.Combobox(self.search_frame, values=["Name", "Gender", "Age", "Breed"])
        self.search_field.grid(row=0, column=1)
        self.search_field.current(0)  # Set default value to first option

        tk.Label(self.search_frame, text="Search Term:").grid(row=1, column=0)
        self.search_entry = tk.Entry(self.search_frame)
        self.search_entry.grid(row=1, column=1)

        tk.Button(self.search_frame, text="Search", command=self.search_pet).grid(row=1, column=2)

        self.results_text = tk.Text(self.search_frame, width=50, height=10)
        self.results_text.grid(row=2, column=0, columnspan=3)

    def browse_photo(self):
        # Define allowed file types
        file_types = [("Image Files", "*.jpg;*.jpeg;*.png"), ("All Files", "*.*")]

        try:
            file_path = filedialog.askopenfilename(filetypes=file_types)
            if file_path:  # Check if a file was selected
                self.pet_photo_entry.delete(0, tk.END)
                self.pet_photo_entry.insert(0, file_path)
            else:
                messagebox.showwarning("File Selection", "No file selected.")
        except Exception as e:
            messagebox.showerror("Error", f"An error occurred: {str(e)}")

    def register_pet(self):
        name = self.pet_name_entry.get()
        gender = self.pet_gender_entry.get()
        age = self.pet_age_entry.get()
        breed = self.pet_breed_entry.get()
        photo = self.pet_photo_entry.get()

        if not all([name, gender, age, breed, photo]):
            messagebox.showwarning("Input Error", "Please fill out all fields.")
            return

        # Create a Pet object
        pet = pet_pb2.RegisterNewPetRequest(
            name=name,
            gender=gender,
            age=int(age),
            breed=breed,
            photo=photo
        )

        # Call gRPC register method
        response = self.stub.RegisterNewPet(pet)
        messagebox.showinfo("Registration", response.message)

    def search_pet(self):
        search_term = self.search_entry.get()
        search_field = self.search_field.get()

        if not search_term:
            messagebox.showwarning("Input Error", "Please enter a search term.")
            return

            # Construct search request based on the selected field
            search_request = None
            if search_field == "Name":
                search_request = pet_pb2.SearchPetRequest(name=search_term)
            elif search_field == "Gender":
                search_request = pet_pb2.SearchPetRequest(gender=search_term)
            elif search_field == "Age":
                search_request = pet_pb2.SearchPetRequest(age=int(search_term))  # Ensure age is an integer
            elif search_field == "Breed":
                search_request = pet_pb2.SearchPetRequest(breed=search_term)

            # Call gRPC search method
            response = self.stub.SearchPet(search_request)

            self.results_text.delete(1.0, tk.END)  # Clear previous results
            if response.pets:
                for pet in response.pets:
                    self.results_text.insert(tk.END,
                                             f"Name: {pet.name}, Gender: {pet.gender}, Age: {pet.age}, Breed: {pet.breed}, Photo: {pet.photo}\n")
            else:
                self.results_text.insert(tk.END, "No pets found.")


if __name__ == "__main__":
    root = tk.Tk()
    app = PetApp(root)
    root.mainloop()
